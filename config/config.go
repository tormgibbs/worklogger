package config

import (
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"github.com/tormgibbs/worklogger/auth"
)

type GithubCreds struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

var (
	Github GithubCreds
	DSN    string
)

func Init() {
	_ = godotenv.Load()

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	clientID := viper.GetString("GITHUB_CLIENT_ID")
	if clientID == "" {
		if id, err := auth.GetToken("github_client_id"); err == nil {
			clientID = id
		}
	}

	clientSecret := viper.GetString("GITHUB_CLIENT_SECRET")
	if clientSecret == "" {
		if secret, err := auth.GetToken("github_client_secret"); err == nil {
			clientSecret = secret
		}
	}

	redirectURI := viper.GetString("GITHUB_REDIRECT_URI")
	if redirectURI == "" {
		if uri, err := auth.GetToken("github_redirect_uri"); err == nil {
			redirectURI = uri
		} else {
			redirectURI = "http://localhost:8080/callback"
		}
	}

	Github = GithubCreds{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURI:  redirectURI,
	}

	DSN = viper.GetString("WORKLOGGER_DSN")
	if DSN == "" {
		DSN = ".worklogger/db.sqlite"
	}
}

// func Init() {
// 	_ = godotenv.Load()

// 	viper.AutomaticEnv()
// 	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

// 	Github = GithubCreds{
// 		ClientID:     viper.GetString("GITHUB_CLIENT_ID"),
// 		ClientSecret: viper.GetString("GITHUB_CLIENT_SECRET"),
// 		RedirectURI:  viper.GetString("GITHUB_REDIRECT_URI"),
// 	}

// 	DSN = viper.GetString("WORKLOGGER_DSN")
// 	if DSN == "" {
// 		DSN = ".worklogger/db.sqlite"
// 	}
// }
