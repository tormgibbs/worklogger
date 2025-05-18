package auth

import "errors"

const (
	GitHubOAuth = "GitHub"
	LocalAuth   = "Local"
	ServiceName = "worklogger"
	LocalKey    = "local_credentials"
	GitHubKey   = "github_token"
)

var (
	ErrInvalidPassword    = errors.New("invalid password")
	ErrUsernameMismatch   = errors.New("username mismatch")
	ErrNoLocalCredentials = errors.New("no local credentials found")
	ErrCredentialsCorrupt = errors.New("stored credentials are corrupted")
	ErrUserAlreadyExists  = errors.New("user already exists")
)
