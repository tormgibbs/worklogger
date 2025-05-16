package auth

import (
	"errors"
	"fmt"

	"github.com/zalando/go-keyring"
	"golang.org/x/crypto/bcrypt"
)

func SaveLocalCredentials(username, password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	creds := fmt.Sprintf("%s|%s", username, hashed)
	return keyring.Set(serviceName, localKey, creds)
}

func VerifyLocalCredentials(username, password string) error {
	creds, err := keyring.Get(serviceName, localKey)
	if err != nil {
		return errors.New("no local credentials found")
	}

	var storedUsername, storedHash string
	fmt.Scanf(creds, "%s|%s", &storedUsername, &storedHash)

	if storedUsername != username {
		return errors.New("username mismatch")
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
	if err != nil {
		return errors.New("invalid password")
	}

	return nil
}
