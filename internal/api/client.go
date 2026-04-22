package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/itshivams/studex-cli/internal/config"
)

const BaseURL = "https://studex.itshivam.in/api"

type UserProfile struct {
	Username       string   `json:"username"`
	FullName       string   `json:"full_name"`
	Role           string   `json:"role"`
	UserPic        string   `json:"userpic"`
	Bio            string   `json:"bio"`
	Github         string   `json:"github"`
	Instagram      string   `json:"instagram"`
	Linkedin       string   `json:"linkedin"`
	Status         string   `json:"status"`
	Location       string   `json:"location"`
	Organization   string   `json:"organization"`
	LastSeen       string   `json:"last_seen"`
	FollowersCount int      `json:"followersCount"`
	FollowingCount int      `json:"followingCount"`
	PostsCount     int      `json:"posts_count"`
	BlogsCount     int      `json:"blogs_count"`
	IsFollowing    bool     `json:"is_following"`
}

type AuthResponse struct {
	Access   string `json:"access"`
	Refresh  string `json:"refresh"`
	Username string `json:"username"`
	ID       string `json:"id"`
}

type LoginPayload struct {
	UsernameOrEmail string `json:"usernameOrEmail"`
	Password        string `json:"password"`
}

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

func Login(usernameOrEmail, password string) (*AuthResponse, error) {
	payload := LoginPayload{
		UsernameOrEmail: usernameOrEmail,
		Password:        password,
	}

	bodyData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/login", BaseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("login failed with status: %d", resp.StatusCode)
	}

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return nil, err
	}

	return &authResp, nil
}

func SearchUser(username string) ([]UserProfile, error) {
	url := fmt.Sprintf("%s/search?username=%s", BaseURL, username)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	token := config.GetToken()
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search failed with status: %d", resp.StatusCode)
	}

	var users []UserProfile
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return nil, err
	}

	return users, nil
}

func SearchUserByName(name string) ([]UserProfile, error) {
	url := fmt.Sprintf("%s/search?name=%s", BaseURL, name)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	token := config.GetToken()
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search failed with status: %d", resp.StatusCode)
	}

	var users []UserProfile
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return nil, err
	}

	return users, nil
}
