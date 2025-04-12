package config

import (
	"errors"
	"fmt"
	"os"
)

var (
	EmptyError    = errors.New("data is empty")
	EmptyKeyError = errors.New("key is empty")
	SecretKey     string
	user          string
	password      string
)

func ConfigProducts() (string, error) {
	user = os.Getenv("user")
	password = os.Getenv("password")
	SecretKey = os.Getenv("secret_key")
	if user == "" || password == "" || SecretKey == "" {
		return "", fmt.Errorf("%w", EmptyError)
	}
	return fmt.Sprintf("host=localhost port=5432 user=%s password=%s dbname=postgres sslmode=disable", user, password), nil
}

func ConfigUser() string {
	return fmt.Sprintf("host=localhost port=5432 user=%s password=%s dbname=users sslmode=disable", user, password)
}
