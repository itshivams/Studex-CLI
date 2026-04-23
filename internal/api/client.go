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

type RegisterStep1Request struct {
	Step         int    `json:"step"`
	Username     string `json:"username"`
	FullName     string `json:"full_name"`
	Email        string `json:"email"`
	PhoneNumber  string `json:"phone_number"`
	ReferralCode string `json:"referralCode"`
}

type RegisterStep1Response struct {
	Message      string `json:"message"`
	RequestId    string `json:"requestId"`
	ResendsLeft  int    `json:"resendsLeft"`
	ExpiresInSec int    `json:"expiresInSec"`
}

type RegisterStep2Request struct {
	Step      int    `json:"step"`
	RequestId string `json:"requestId"`
	Otp       string `json:"otp"`
}

type RegisterStep2Response struct {
	Message   string `json:"message"`
	TempToken string `json:"tempToken"`
}

type RegisterStep3Request struct {
	Step      int    `json:"step"`
	TempToken string `json:"tempToken"`
	Password  string `json:"password"`
}

type RegisterStep3Response struct {
	Message string `json:"message"`
	Token   string `json:"token"`
	User    struct {
		Username string `json:"username"`
	} `json:"user"`
}

func doRegisterCall(payload interface{}, response interface{}) error {
	bodyData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/register", BaseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var errResp struct {
			Message string `json:"message"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err == nil && errResp.Message != "" {
			return fmt.Errorf("%s", errResp.Message)
		}
		return fmt.Errorf("request failed with status: %d", resp.StatusCode)
	}

	return json.NewDecoder(resp.Body).Decode(response)
}

func RegisterStep1(req RegisterStep1Request) (*RegisterStep1Response, error) {
	req.Step = 1
	var res RegisterStep1Response
	if err := doRegisterCall(req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func RegisterStep2(req RegisterStep2Request) (*RegisterStep2Response, error) {
	req.Step = 2
	var res RegisterStep2Response
	if err := doRegisterCall(req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func RegisterStep3(req RegisterStep3Request) (*RegisterStep3Response, error) {
	req.Step = 3
	var res RegisterStep3Response
	if err := doRegisterCall(req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func CheckUsername(username string) (bool, string, error) {
	url := fmt.Sprintf("%s/check-username?username=%s", BaseURL, username)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, "", err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return false, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, "", fmt.Errorf("check username failed with status: %d", resp.StatusCode)
	}

	var res struct {
		Available bool   `json:"available"`
		Message   string `json:"message"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return false, "", err
	}

	return res.Available, res.Message, nil
}
