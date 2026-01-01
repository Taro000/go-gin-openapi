package repo

import (
	"context"
	"database/sql"
)

type GoodluckRepo struct {
	db *sql.DB
}

func (r *GoodluckRepo) Create(ctx context.Context, userID, todoID string) error {
	_, err := r.db.ExecContext(ctx, `INSERT IGNORE INTO goodlucks (user, todo) VALUES (?, ?)`, userID, todoID)
	return err
}

func (r *GoodluckRepo) Delete(ctx context.Context, userID, todoID string) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM goodlucks WHERE user = ? AND todo = ?`, userID, todoID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}


