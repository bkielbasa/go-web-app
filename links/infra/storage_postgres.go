package infra

import (
	"context"
	"database/sql"
	"fmt"
)

type postgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(db *sql.DB) postgresStorage {
	return postgresStorage{
		db: db,
	}
}

func (storage postgresStorage) Add(ctx context.Context, targetLink, id string, tags []string) error {
	q := `INSERT INTO links (id, target_link) values ($1, $2)`
	_, err := storage.db.ExecContext(ctx, q, id, targetLink)
	if err != nil {
		return fmt.Errorf("cannot insert new link to the DB: %w", err)
	}

	return nil
}

func (storage postgresStorage) Get(ctx context.Context, id string) (string, error) {
	q := `SELECT target_link FROM links WHERE id = $1`
	var target string

	err := storage.db.QueryRowContext(ctx, q, id).Scan(&target)
	if err != nil {
		return "", fmt.Errorf("cannot get the result: %w", err)
	}

	return target, nil
}
