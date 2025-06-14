package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/browser"
)

func StartGitHubOAuth(clientID, clientSecret, redirectURI string) error {

	authUrl := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=read:user",
		clientID, url.QueryEscape(redirectURI),
	)

	fmt.Println("\n🌐 Opening GitHub OAuth page...")
	err := browser.OpenURL(authUrl)
	if err != nil {
		return err
	}

	server := &http.Server{Addr: ":3000"}

	done := make(chan bool)

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "No code found", http.StatusBadRequest)
			return
		}

		token, err := exchangeCodeForToken(code, clientID, clientSecret, redirectURI)
		if err != nil {
			fmt.Println("Error exchanging token:", err)
			http.Error(w, "OAuth failed", http.StatusInternalServerError)
			return
		}

		SetToken(GitHubKey, token)
		SaveSession(Session{
			Method:        "github",
			Authenticated: true,
		})

		fmt.Fprintf(w, "You’re logged in! You can close this tab.")

		// Signal done
		go func() {
			done <- true
		}()

	})

	fmt.Println("🚪 Waiting for GitHub OAuth callback on :3000...")

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("HTTP server error:", err)
		}
	}()

	<-done

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	return nil
}

func exchangeCodeForToken(code, clientID, clientSecret, redirectURI string) (string, error) {

	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("code", code)
	data.Set("redirect_uri", redirectURI)

	request, err := http.NewRequest(http.MethodPost, "https://github.com/login/oauth/access_token", strings.NewReader(data.Encode()))

	if err != nil {
		return "", err
	}

	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := http.DefaultClient
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, _ := io.ReadAll(response.Body)
	result := make(map[string]any)
	json.Unmarshal(body, &result)

	token, ok := result["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("no access token returned")
	}

	return token, nil
}
