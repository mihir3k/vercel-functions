package handler

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type authResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
	ExpiresIn   int    `json:"expires_in"`
}

type currentlyPlaying struct {
	IsPlaying            bool   `json:"is_playing"`
	Item                 track  `json:"item"`
	CurrentlyPlayingType string `json:"currently_playing_type"`
}

type track struct {
	Name        string      `json:"name"`
	Album       album       `json:"album"`
	Artists     []artist    `json:"artists"`
	ExternalURL externalURL `json:"external_urls"`
}

type album struct {
	Images []image `json:"images"`
}

type image struct {
	Height int    `json:"height"`
	Weight int    `json:"weight"`
	URL    string `json:"url"`
}

type artist struct {
	Name string `json:"name"`
}

type externalURL struct {
	Spotify string `json:"spotify"`
}

const authURL = "https://accounts.spotify.com/api/token"
const currentlyPlayingURL = "https://api.spotify.com/v1/me/player/currently-playing"

func getAccessToken() (string, error) {
	spotifyClientID := os.Getenv("SPOTIFY_CLIENT_ID")
	spotifyClientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")
	refreshToken := os.Getenv("SPOTIFY_REFRESH_TOKEN")

	authCode := base64.StdEncoding.EncodeToString(([]byte(spotifyClientID + ":" + spotifyClientSecret)))
	authHeader := fmt.Sprintf("Basic %s", authCode)
	params := url.Values{}
	params.Set("grant_type", "refresh_token")
	params.Set("refresh_token", refreshToken)

	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest(http.MethodPost, authURL, strings.NewReader(params.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", authHeader)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	authResponseObj := &authResponse{}

	if res.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return "", err
		}
		err = json.Unmarshal(body, authResponseObj)
		if err != nil {
			return "", err
		}
	}
	return authResponseObj.AccessToken, nil
}

func getCurrentlyPlaying() (*currentlyPlaying, error) {
	accessToken, err := getAccessToken()
	if err != nil {
		return nil, err
	}
	authHeader := fmt.Sprintf("Bearer %s", accessToken)

	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest(http.MethodGet, currentlyPlayingURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", authHeader)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	currentlyPlayingObj := &currentlyPlaying{}

	if res.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(body, currentlyPlayingObj)
		if err != nil {
			return nil, err
		}
	}
	return currentlyPlayingObj, nil
}

// Handler godoc
func Handler(w http.ResponseWriter, r *http.Request) {
	currentlyPlaying, err := getCurrentlyPlaying()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(currentlyPlaying)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
