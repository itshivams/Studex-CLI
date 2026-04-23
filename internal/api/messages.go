package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/itshivams/studex-cli/internal/config"
)

type EncryptedData struct {
	Ciphertext   string `json:"ciphertext"`
	Iv           string `json:"iv"`
	RecipientKey string `json:"recipientKey"`
	SenderKey    string `json:"senderKey"`
}

type Message struct {
	ID            string         `json:"_id"`
	From          string         `json:"from"`
	To            string         `json:"to"`
	Text          string         `json:"text"`
	EncryptedData *EncryptedData `json:"encryptedData"`
	Type          string         `json:"type"`
	CreatedAt     string         `json:"createdAt"`
	Read          bool           `json:"read"`
	FromID        string         `json:"fromId"`
	ToID          string         `json:"toId"`
}

type ThreadResponse struct {
	Messages []Message       `json:"messages"`
	Shares   []interface{}   `json:"shares"`
}

func GetMessageThread(withUsername string) (*ThreadResponse, error) {
	url := fmt.Sprintf("%s/messages/thread?with=%s&include=shares", BaseURL, withUsername)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	token := config.GetToken()
	if token == "" {
		return nil, fmt.Errorf("no authentication token found")
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch thread with status: %d", resp.StatusCode)
	}

	var thread ThreadResponse
	if err := json.NewDecoder(resp.Body).Decode(&thread); err != nil {
		return nil, err
	}

	return &thread, nil
}

type SendMessagePayload struct {
	ToUsername string `json:"toUsername"`
	Text       string `json:"text"`
}

type SendMessageResponse struct {
	Message string  `json:"message"`
	Data    Message `json:"data"`
}

func SendMessage(toUsername, text string) (*SendMessageResponse, error) {
	payload := SendMessagePayload{
		ToUsername: toUsername,
		Text:       text,
	}

	bodyData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/messages/send", BaseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	token := config.GetToken()
	if token == "" {
		return nil, fmt.Errorf("no authentication token found")
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to send message with status: %d", resp.StatusCode)
	}

	var res SendMessageResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}
