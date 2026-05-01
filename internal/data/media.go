package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type MediaModel struct {
	DB *sql.DB
}

type Media struct {
	ID        uuid.UUID    `json:"-"`
	PublicID  string       `json:"public_id"`
	UserID    uuid.UUID    `json:"-"`
	Filename  string       `json:"filename"`
	MimeType  string       `json:"mime_type"`
	Size      int64        `json:"size"`
	CreatedAt sql.NullTime `json:"created_at"`
	Version   int32        `json:"-"`
}

func (m MediaModel) CreateMedia(media *Media) error {
	query := `
		INSERT INTO media (public_media_id, user_id, filename, mime_type, size, created_at, version)
		VALUES ($1, $2, $3, $4, $5, $6, 1)
		RETURNING id, created_at, version
	`

	args := []any{media.PublicID, media.UserID, media.Filename, media.MimeType, media.Size, media.CreatedAt}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&media.ID, &media.CreatedAt, &media.Version)

}

func (m MediaModel) GetMediaByPublicID(public_id string) (*Media, error) {
	query := `
		SELECT id, public_media_id, user_id, filename, mime_type, size, created_at, version
		FROM media
		WHERE public_media_id = $1
	`

	var media Media

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, public_id).Scan(
		&media.ID,
		&media.PublicID,
		&media.UserID,
		&media.Filename,
		&media.MimeType,
		&media.Size,
		&media.CreatedAt,
		&media.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &media, nil
}

func (m MediaModel) UpdateMedia(media *Media) error {

	query := `
		UPDATE media
		SET filename = $1, mime_type = $2, size = $3, version = version + 1
		WHERE public_media_id = $4 AND version = $5
		RETURNING version
	`

	args := []any{media.Filename, media.MimeType, media.Size, media.PublicID, media.Version}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&media.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (m MediaModel) GetMedia(id uuid.UUID) (*Media, error) {
	query := `
		SELECT id, public_media_id, user_id, filename, mime_type, size, created_at, version
		FROM media
		WHERE id = $1
	`

	var media Media

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&media.ID,
		&media.PublicID,
		&media.UserID,
		&media.Filename,
		&media.MimeType,
		&media.Size,
		&media.CreatedAt,
		&media.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &media, nil

}

func (m MediaModel) ListMediaByUser(id string) ([]*Media, error) {
	query := `
		SELECT m.id, m.public_media_id, m.user_id, m.filename, m.mime_type, m.size, m.created_at, m.version
		FROM media m
		JOIN users ON m.user_id = users.id
		WHERE users.public_id = $1
		ORDER BY m.created_at DESC
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	medium := []*Media{}
	for rows.Next() {
		var media Media
		err := rows.Scan(
			&media.ID,
			&media.PublicID,
			&media.UserID,
			&media.Filename,
			&media.MimeType,
			&media.Size,
			&media.CreatedAt,
			&media.Version,
		)
		if err != nil {
			return nil, err
		}
		medium = append(medium, &media)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return medium, nil
}

func (m MediaModel) DeleteMedia(id uuid.UUID) error {
	query := `
		DELETE FROM media
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (m MediaModel) GetAll(title string, filetype string, filters Filters) ([]*Media, error) {
	query := fmt.Sprintf(`
		SELECT id, public_media_id, user_id, filename, mime_type, size, created_at, version
		FROM media
		WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (mime_type = $2 OR $2 = '')
		ORDER BY %s %s`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	medium := []*Media{}

	for rows.Next() {
		var media Media

		err := rows.Scan(
			&media.ID,
			&media.PublicID,
			&media.UserID,
			&media.Filename,
			&media.MimeType,
			&media.Size,
			&media.CreatedAt,
			&media.Version,
		)
		if err != nil {
			return nil, err
		}
		medium = append(medium, &media)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return medium, nil
}
