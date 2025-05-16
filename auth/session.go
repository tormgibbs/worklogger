package auth

import (
	"encoding/json"
	"os"
	"path/filepath"
)


type Session struct {
	Method        string `json:"method"`
	Authenticated bool   `json:"authenticated"`
}

func getSessionFilePath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	dir := filepath.Join(configDir, "worklogger")
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", err
	}

	return filepath.Join(dir, "auth.json"), nil
}

func SaveSession(s Session) error {
	path, err := getSessionFilePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(s, "", " ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}

func LoadSession() (Session, error) {
	s := Session{}

	path, err := getSessionFilePath()
	if err != nil {
		return s, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return s, nil
	}

	err = json.Unmarshal(data, &s)

	return s, err
}

func DeleteSession() error {
	path, err := getSessionFilePath()
	if err != nil {
		return err
	}
	return os.Remove(path)
}
