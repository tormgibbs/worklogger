package auth

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/x/term"
	"github.com/zalando/go-keyring"
	"golang.org/x/crypto/bcrypt"
)

func SaveLocalCredentials(username, password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	creds := fmt.Sprintf("%s|%s", username, hashed)
	return keyring.Set(ServiceName, LocalKey, creds)
}

func VerifyLocalCredentials(username, password string) error {
	creds, err := keyring.Get(ServiceName, LocalKey)
	if err != nil {
		return errors.New("no local credentials found")
	}

	parts := strings.SplitN(creds, "|", 2)
	if len(parts) != 2 {
		return errors.New("stored credentials are corrupted")
	}

	storedUsername := parts[0]
	storedHash := parts[1]

	if storedUsername != username {
		return errors.New("username mismatch")
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
	if err != nil {
		return errors.New("invalid password")
	}

	return nil
}

func LocalLogin() error {
	username, err := promptUsername()
	if err != nil {
		return err
	}

	password, err := promptPassword()
	if err != nil {
		return err
	}

	if err := VerifyLocalCredentials(username, password); err != nil {
		return fmt.Errorf("login failed: %w", err)
	}

	if err := startSession(); err != nil {
		return fmt.Errorf("failed to save session: %w", err)
	}

	fmt.Println("Logged in successfully with local credentials.")
	return nil
}

func LocalSignUp() error {
	username, err := promptUsername()
	if err != nil {
		return err
	}

	password, err := promptPassword()
	if err != nil {
		return err
	}

	if err := SaveLocalCredentials(username, password); err != nil {
		return fmt.Errorf("signing up failed: %w", err)
	}

	if err := startSession(); err != nil {
		return fmt.Errorf("failed to save session: %w", err)
	}

	fmt.Println("Signup successful. You're now logged in locally.")
	return nil
}

func promptUsername() (string, error) {
	fmt.Print("Enter your username: ")
	var username string
	if _, err := fmt.Scanln(&username); err != nil {
		return "", fmt.Errorf("failed to read username: %w", err)
	}
	return strings.TrimSpace(username), nil
}

func promptPassword() (string, error) {
	fmt.Print("Enter your password: ")
	bytePass, err := term.ReadPassword(uintptr(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		return "", fmt.Errorf("failed to read password: %w", err)
	}
	return string(bytePass), nil
}

func startSession() error {
	return SaveSession(Session{
		Method:        LocalAuth,
		Authenticated: true,
	})
}
