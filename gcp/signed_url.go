package gcp

import (
	"fmt"
	"io/ioutil"
	"signdoc_api/config"
	"time"

	"cloud.google.com/go/storage"
	"golang.org/x/oauth2/google"
)

// https://cloud.google.com/storage/docs/access-control/signed-urls

// GenerateV4GetObjectSignedURL generates object signed URL with GET method.
func GenerateV4GetObjectSignedURL(bucket, object string) (string, error) {
	jsonKey, err := ioutil.ReadFile(config.C.GCP.GoogleApplicationCredentials)
	if err != nil {
		return "", fmt.Errorf("ioutil.ReadFile: %v", err)
	}
	conf, err := google.JWTConfigFromJSON(jsonKey)
	if err != nil {
		return "", fmt.Errorf("google.JWTConfigFromJSON: %v", err)
	}
	opts := &storage.SignedURLOptions{
		Scheme:         storage.SigningSchemeV4,
		Method:         "GET",
		GoogleAccessID: conf.Email,
		PrivateKey:     conf.PrivateKey,
		Expires:        time.Now().Add(15 * time.Minute),
	}
	u, err := storage.SignedURL(bucket, object, opts)
	if err != nil {
		return "", fmt.Errorf("storage.SignedURL: %v", err)
	}

	return u, nil
}
