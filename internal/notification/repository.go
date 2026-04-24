package notification

import (
	"database/sql"
	"encoding/json"
	"log"

	"ProjectNotification/internal/models"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Save(n models.Notification) error {
	payloadJSON, err := json.Marshal(n.Payload)
	if err != nil {
		return err
	}

	query := `
	INSERT INTO notification (id, user_id, type, payload, status, created_at)
	VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err = r.db.Exec(query, n.ID, n.UserID, n.Type, payloadJSON, n.Status, n.CreatedAt)
	if err != nil {
		log.Println("failed to insert notification:", err)
		return err
	}

	log.Println("notification saved:", n.ID)
	return nil
}

func (r *Repository) UpdateStatus(id string, status string) error {
	_, err := r.db.Exec(`UPDATE notification SET status = $2 WHERE id = $1`, id, status)
	if err != nil {
		log.Println("failed to update notification status:", err)
		return err
	}
	return nil
}
