package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Post struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	UserID    int64     `json:"user_id"`
	Tags      []string  `json:"tags"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	Version   int       `json:"version"`
	Comments  []Comment `json:"comments"`
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

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var createdAt, updatedAt pgtype.Timestamptz
	err := s.db.QueryRow(
		ctx,
		query,
		post.Title,
		post.Content,
		post.UserID,
		post.Tags,
		// pgtype.Array[string]{Elements: post.Tags},
	).Scan(&post.ID, &createdAt, &updatedAt)
	// Scan：把資料庫返回的數據讀取到 Go 變量中
	// &：告訴 Scan 要把數據寫到哪個變量的地址
	// Scan 的參數順序必須和 SQL RETURNING 的順序一致

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

func (s *PostsStore) GetByID(ctx context.Context, postID int64) (*Post, error) {
	query := `
		SELECT id, title, content, user_id, tags, created_at, updated_at, version
		FROM posts
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var createdAt, updatedAt pgtype.Timestamptz
	var post Post
	err := s.db.QueryRow(
		ctx,
		query,
		postID,
	).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.UserID,
		&post.Tags,
		&createdAt,
		&updatedAt,
		&post.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	if createdAt.Valid {
		post.CreatedAt = createdAt.Time.Format(time.RFC3339)
	}
	if updatedAt.Valid {
		post.UpdatedAt = updatedAt.Time.Format(time.RFC3339)
	}

	return &post, nil
}

func (s *PostsStore) Delete(ctx context.Context, postID int64) error {
	query := `DELETE FROM posts WHERE id = $1`

	// 為資料庫操作設定一個 5 秒的超時限制
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res, err := s.db.Exec(ctx, query, postID)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		fmt.Println("測試: ", res.RowsAffected())
		return ErrNotFound
	}

	return nil
}

// 它通過 修改傳入的指針指向的對象 來更新數據
func (s *PostsStore) Update(ctx context.Context, post *Post) error {
	query := `
		UPDATE posts
		SET title = $1, content = $2, version = version + 1
		WHERE id = $3 AND version = $4
		RETURNING version
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRow(
		ctx,
		query,
		post.Title,
		post.Content,
		post.ID,
		post.Version,
	).Scan(&post.Version)
	if err != nil {

		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
	}

	return nil
}
