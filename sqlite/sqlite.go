package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/empfaze/golang_bot/lib/storage"
)

type Storage struct {
	db *sql.DB
}

func New(path string) (*Storage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("Couldn't open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("Couldn't connect to database: %w", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Save(ctx context.Context, p *storage.Page) error {
	query := `INSERT INTO pages (url, user_name) VALUES (?, ?)`

	if _, err := s.db.ExecContext(ctx, query, p.URL, p.Username); err != nil {
		return fmt.Errorf("Couldn't save page: %w", err)
	}

	return nil
}

func (s *Storage) PickRandom(ctx context.Context, username string) (*storage.Page, error) {
	query := `SELECT url FROM pages WHERE user_name = ? ORDER BY RANDOM() LIMIT 1`

	var url string

	err := s.db.QueryRowContext(ctx, query, username).Scan(&url)
	if err == sql.ErrNoRows {
		return nil, storage.ErrNoSavedPages
	}
	if err != nil {
		return nil, fmt.Errorf("Couldn't pick random page: %w", err)
	}

	return &storage.Page{URL: url, Username: username}, nil
}

func (s *Storage) Remove(ctx context.Context, page *storage.Page) error {
	query := `DELETE FROM pages WHERE url = ? AND user_name = ?`
	if _, err := s.db.ExecContext(ctx, query, page.URL, page.Username); err != nil {
		return fmt.Errorf("Couldn't remove page: %w", err)
	}

	return nil
}

func (s *Storage) IsExist(ctx context.Context, page *storage.Page) (bool, error) {
	query := `SELECT COUNT(*) FROM pages WHERE url = ? AND user_name = ?`

	var count int

	if err := s.db.QueryRowContext(ctx, query, page.URL, page.Username).Scan(&count); err != nil {
		return false, fmt.Errorf("Couldn't check if page exists: %w", err)
	}

	return count > 0, nil
}

func (s *Storage) Init(ctx context.Context) error {
	query := `CREATE TABLE IF NOT EXISTS pages (url TEXT, user_name TEXT)`

	if _, err := s.db.ExecContext(ctx, query); err != nil {
		return fmt.Errorf("Couldn't create table: %w", err)
	}

	return nil
}
