package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Post struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	UserID    int64    `json:"user_id"`
	Tags      []string `json:"tags"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

type PostsStore struct {
	db *pgxpool.Pool
}

func (s *PostsStore) Create(ctx context.Context, post *Post) error {
	query := `
		INSERT INTO posts (title, content, user_id, tags)
		VALUES ($1, $2, $3, $4)
		RETURNING id,  created_at, updated_at
	`

	var createdAt, updatedAt pgtype.Timestamptz
	err := s.db.QueryRow(
		ctx,
		query,
		post.Title,
		post.Content,
		post.UserID,
		pgtype.Array[string]{Elements: post.Tags},
	).Scan(&post.ID, &createdAt, &updatedAt)

	if err != nil {
		return err
	}

	if createdAt.Valid {
		post.CreatedAt = createdAt.Time.Format(time.RFC3339)
	}
	if updatedAt.Valid {
		post.UpdatedAt = updatedAt.Time.Format(time.RFC3339)
	}

	return nil
}
