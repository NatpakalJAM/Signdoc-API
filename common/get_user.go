package common

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func GetUser(c *fiber.Ctx) (user string) {
	user = c.Get("username")
	temp := strings.Split(user, "@")
	user = temp[0]
	return user
}
