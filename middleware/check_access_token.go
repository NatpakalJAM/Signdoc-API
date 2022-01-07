package middleware

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"signdoc_api/config"
	"signdoc_api/custom_cache"
	"signdoc_api/model"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

var cacheDuration = time.Minute * 800

type responseValidateToken struct {
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func CheckAccessToken(c *fiber.Ctx) error {
	// skip CheckAccessToken when development
	if config.C.Environment == "development" {
		return c.Next()
	}

	username := c.Get("username")
	accessToken := c.Get("access_token")
	if username == "" || accessToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.JSONResponse{
			Code:    fiber.StatusBadRequest,
			Success: false,
			Message: "Header `access_token` and `username` are required.",
			Data:    nil,
		})
	}

	// check cache auth validate
	dataCache, found := custom_cache.Ca.Get(username)
	if found && accessToken == dataCache {
		return c.Next()
	}

	form := url.Values{
		"username":     {username},
		"access_token": {accessToken},
	}
	client := &http.Client{}
	URL := fmt.Sprintf("%s/oauth2/validate", config.C.Auth.URL)
	req, err := http.NewRequest("POST", URL, strings.NewReader(form.Encode()))
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", config.C.Auth.Authorization))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.JSONResponse{
			Code:    fiber.StatusInternalServerError,
			Success: false,
			Message: fmt.Sprintf("Error CheckAccessToken: NewRequest => %v", err),
			Data:    nil,
		})
	}
	resp, err := client.Do(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.JSONResponse{
			Code:    fiber.StatusInternalServerError,
			Success: false,
			Message: fmt.Sprintf("Error CheckAccessToken: Do Request => %v", err),
			Data:    nil,
		})
	}
	defer resp.Body.Close()

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.JSONResponse{
			Code:    fiber.StatusInternalServerError,
			Success: false,
			Message: fmt.Sprintf("Error CheckAccessToken: ioutil.ReadAll => %v", err),
			Data:    nil,
		})
	}
	var respValidateToken responseValidateToken
	_ = json.Unmarshal(response, &respValidateToken)

	if !respValidateToken.Success {
		return c.Status(fiber.StatusUnauthorized).JSON(model.JSONResponse{
			Code:    fiber.StatusUnauthorized,
			Success: false,
			Message: "Error CheckAccessToken: Unauthorized",
			Data:    nil,
		})
	}

	custom_cache.Ca.Set(username, accessToken, cacheDuration)

	return c.Next()
}
