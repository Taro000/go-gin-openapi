package repo

import (
	"context"
	"database/sql"
)

type TodoStatusRepo struct {
	db *sql.DB
}

func (r *TodoStatusRepo) Ensure(ctx context.Context, status string) error {
	// Allow any status code (OpenAPI doesn't define enum). This keeps FK valid.
	_, err := r.db.ExecContext(ctx, `INSERT IGNORE INTO todo_statuses (status) VALUES (?)`, status)
	return err
}


