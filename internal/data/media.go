package data

import (
	"database/sql"

	"github.com/google/uuid"
)

type MediaModel struct {
	DB *sql.DB
}

type Media struct {
	ID        uuid.UUID    `json:"-"`
	PublicID  string       `json:"public_id"`
	UserID    uuid.UUID    `json:"user_id"`
	Filename  string       `json:"filename"`
	MimeType  string       `json:"mime_type"`
	Size      int64        `json:"size"`
	CreatedAt sql.NullTime `json:"created_at"`
	Version   int32        `json:"-"`
}

func (m MediaModel) CreateMedia(media *Media) error {
	return nil
}

func (m MediaModel) GetMediaNameByPublicID(id string) (string, error) {
	return "", nil
}

func (m MediaModel) ListMediaByUser(id string) ([]*Media, error) {
	return nil, nil
}
