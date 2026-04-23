package api

import (
	"encoding/json"
	"fmt"
	"time"

	"golang.org/x/net/websocket"
)

const DiscussionWSURL = "wss://chat.studex.itshivam.in/ws"
const DiscussionOrigin = "https://studex.itshivam.in"

type ChatMessage struct {
	Username  string `json:"username"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
	Userpic   string `json:"userpic,omitempty"`
	Verified  bool   `json:"verified,omitempty"`
	Type      string `json:"type,omitempty"` 
	ActiveUsers int `json:"activeUsers,omitempty"`
}

type DiscussionConn struct {
	ws *websocket.Conn
}

func ConnectDiscussion() (*DiscussionConn, error) {
	cfg, err := websocket.NewConfig(DiscussionWSURL, DiscussionOrigin)
	if err != nil {
		return nil, fmt.Errorf("ws config error: %w", err)
	}

	ws, err := websocket.DialConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("ws dial error: %w", err)
	}
	return &DiscussionConn{ws: ws}, nil
}

func (d *DiscussionConn) Close() {
	if d.ws != nil {
		d.ws.Close()
	}
}

func (d *DiscussionConn) Send(msg ChatMessage) error {
	if msg.Timestamp == "" {
		msg.Timestamp = time.Now().Format("03:04 PM")
	}
	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	_, err = d.ws.Write(b)
	return err
}

func (d *DiscussionConn) Receive() (*ChatMessage, error) {
	var buf [65536]byte
	n, err := d.ws.Read(buf[:])
	if err != nil {
		return nil, err
	}
	var msg ChatMessage
	if err := json.Unmarshal(buf[:n], &msg); err != nil {
		msg = ChatMessage{Content: string(buf[:n])}
	}
	return &msg, nil
}
