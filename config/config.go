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

	Github = GithubCreds{
		ClientID:     clientID,
		ClientSecret: clientSecret,
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
// 	}

// 	DSN = viper.GetString("WORKLOGGER_DSN")
// 	if DSN == "" {
// 		DSN = ".worklogger/db.sqlite"
// 	}

// 	fmt.Println("clientID:", Github.ClientID)
// 	fmt.Println("clientSecret:", Github.ClientSecret)
// 	fmt.Println("viper-clientID:", viper.GetString("GITHUB_CLIENT_ID"))
// 	fmt.Println("viper-clientSecret:", viper.GetString("GITHUB_CLIENT_SECRET"))
// }
