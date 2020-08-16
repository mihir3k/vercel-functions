package handler

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type authResponse struct {
	accessToken string
	tokenType   string
	scope       string
	expiresIn   int
}

func getAccessToken(w http.ResponseWriter, r *http.Request) error {
	spotifyClientID := os.Getenv("SPOTIFY_CLIENT_ID")
	spotifyClientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")
	refreshToken := os.Getenv("SPOTIFY_REFRESH_TOKEN")

	apiURL := "https://accounts.spotify.com/api/token"

	authCode := base64.StdEncoding.EncodeToString(([]byte(spotifyClientID + ":" + spotifyClientSecret)))
	authHeader := fmt.Sprintf("Basic %s", authCode)

	params := url.Values{}
	params.Set("grant_type", "refresh_token")
	params.Set("refresh_token", refreshToken)

	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest(http.MethodPost, apiURL, strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", authHeader)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		bodyString := string(body)
		fmt.Fprintf(w, bodyString)
	}
	return nil
}

// Handler returns a string
func Handler(w http.ResponseWriter, r *http.Request) {
	err := getAccessToken(w, r)
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
}
