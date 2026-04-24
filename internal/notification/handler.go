package notification

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"ProjectNotification/internal/delivery/telegram"
	"ProjectNotification/internal/models"
	"ProjectNotification/internal/redis"

	"github.com/google/uuid"
)

var ErrUnsupportedEventType = errors.New("unsupported event type")

type Handler struct {
	repo    *Repository
	redis   *redis.Client
	telegram *telegram.Sender
}

func NewHandler(repo *Repository, redisClient *redis.Client, telegramSender *telegram.Sender) *Handler {
	return &Handler{repo: repo, redis: redisClient, telegram: telegramSender}
}

func (h *Handler) Handle(ctx context.Context, event models.Event) error {
	// Map event -> notification by type.
	switch event.EventType {
	case "user_registered":
		// supported
	default:
		return ErrUnsupportedEventType
	}

	// Dedupe based on stable hash of event content.
	if h.redis != nil {
		sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%d|%v", event.EventType, event.UserID, event.Timestamp, event.Payload)))
		key := "dedupe:event:" + hex.EncodeToString(sum[:])
		ok, err := h.redis.SetNX(ctx, key, "1", 24*time.Hour)
		if err != nil {
			// Redis is optional: if it's down, keep processing without dedupe.
			ok = true
		}
		if !ok {
			return nil // already processed
		}
	}

	n := models.Notification{
		ID:        uuid.NewString(),
		UserID:    event.UserID,
		Type:      event.EventType,
		Payload:   event.Payload,
		Status:    "pending",
		CreatedAt: time.Now().Unix(),
	}

	if err := h.repo.Save(n); err != nil {
		return err
	}

	// Deliver via Telegram (best-effort if not configured).
	if h.telegram != nil {
		text := fmt.Sprintf("New notification\nType: %s\nUser: %s\nPayload: %v", n.Type, n.UserID, n.Payload)
		if err := h.telegram.SendMessage(ctx, text); err != nil {
			_ = h.repo.UpdateStatus(n.ID, "failed")
			return err
		}
		_ = h.repo.UpdateStatus(n.ID, "sent")
	}

	return nil
}
