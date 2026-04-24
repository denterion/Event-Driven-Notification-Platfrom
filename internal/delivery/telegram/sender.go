package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
)

var ErrTelegramNotConfigured = errors.New("telegram is not configured (missing TELEGRAM_BOT_TOKEN or TELEGRAM_CHAT_ID)")

type Sender struct {
	token  string
	chatID string
	client *http.Client
}

func NewFromEnv() *Sender {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	chatID := os.Getenv("TELEGRAM_CHAT_ID")

	return &Sender{
		token:  token,
		chatID: chatID,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *Sender) IsConfigured() bool {
	return s.token != "" && s.chatID != ""
}

func (s *Sender) SendMessage(ctx context.Context, text string) error {
	if !s.IsConfigured() {
		return ErrTelegramNotConfigured
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", s.token)
	body, err := json.Marshal(map[string]string{
		"chat_id": s.chatID,
		"text":    text,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		// Avoid leaking token via error strings that might include URL.
		return fmt.Errorf("telegram sendMessage request failed: %T", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("telegram sendMessage failed: status=%s", resp.Status)
	}
	return nil
}

