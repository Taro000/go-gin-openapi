package repo

import (
	"context"
	"database/sql"
)

type TodoStatusRepo struct {
	db *sql.DB
}

func (r *TodoStatusRepo) Ensure(ctx context.Context, status string) error {
	// Ensure status code exists in todo_statuses (keeps FK valid).
	_, err := r.db.ExecContext(ctx, `INSERT IGNORE INTO todo_statuses (status) VALUES (?)`, status)
	return err
}


