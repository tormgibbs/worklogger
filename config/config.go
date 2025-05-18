package config

import (
	"log"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type GithubCreds struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

var Github GithubCreds

func Init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system env vars")
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	Github = GithubCreds{
		ClientID:     viper.GetString("GITHUB_CLIENT_ID"),
		ClientSecret: viper.GetString("GITHUB_CLIENT_SECRET"),
		RedirectURI:  viper.GetString("GITHUB_REDIRECT_URI"),
	}
}
