package auth

import (
	"fmt"

	"github.com/zalando/go-keyring"
)


func SetToken(key string, token string) error {
	return keyring.Set(ServiceName, key, token)
}

func GetToken(key string) (string, error) {
	token, err := keyring.Get(ServiceName, key)
	if err != nil {
		return "", fmt.Errorf("could not retrieve token: %w", err)
	}
	return token, nil
}

func DeleteToken(key string) error {
	return keyring.Delete(ServiceName, key)
}


// var (
// 	hashFlag, messageFlag, authorFlag, dateFLag string
// )

// readCommitCmd.Flags().StringVarP(&hashFlag, "hash", "h", "", "Git commit hash")
// 	readCommitCmd.Flags().StringVarP(&messageFlag, "message", "m", "", "Git commit message")
// 	readCommitCmd.Flags().StringVarP(&authorFlag, "author", "a", "", "Commit author")
// 	readCommitCmd.Flags().StringVarP(&hashFlag, "date", "d", "", "Commit date")

// 	if hashFlag == "" || messageFlag == "" || authorFlag == "" || dateFLag == "" {
// 		fmt.Println("Missing required commit fields")
// 		flag.Usage()
// 		return
