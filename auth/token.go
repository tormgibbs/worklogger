package auth

import (
	"fmt"

	"github.com/zalando/go-keyring"
)


func SetToken(key string, token string) error {
	return keyring.Set(serviceName, key, token)
}

func GetToken(key string) (string, error) {
	token, err := keyring.Get(serviceName, key)
	if err != nil {
		return "", fmt.Errorf("could not retrieve token: %w", err)
	}
	return token, nil
}

func DeleteToken(key string) error {
	return keyring.Delete(serviceName, key)
}
