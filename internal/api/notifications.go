package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/itshivams/studex-cli/internal/config"
)

type NotificationTarget struct {
	Type   string `json:"type"`
	UserID string `json:"user_id"`
}

type NotificationSender struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type NotificationItem struct {
	ID              string             `json:"id"`
	Title           string             `json:"title"`
	Message         string             `json:"message"`
	Target          NotificationTarget `json:"target"`
	Sender          NotificationSender `json:"sender"`
	RecipientsCount int                `json:"recipients_count"`
	Status          string             `json:"status"`
	CreatedAt       string             `json:"created_at"`
	SentAt          string             `json:"sent_at"`
	Read            bool               `json:"read"`
}

type NotificationsResponse struct {
	Items []NotificationItem `json:"items"`
}

func GetNotifications() ([]NotificationItem, error) {
	url := fmt.Sprintf("%s/notifications", BaseURL)

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
		return nil, fmt.Errorf("failed to fetch notifications with status: %d", resp.StatusCode)
	}

	var res NotificationsResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return res.Items, nil
}
