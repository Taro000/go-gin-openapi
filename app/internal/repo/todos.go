package repo

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Todo struct {
	ID          string
	Owner       string
	Status      string
	Title       string
	Content     string
	DueDatetime *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type TodoRepo struct {
	db *sql.DB
}

func (r *TodoRepo) Create(ctx context.Context, t Todo) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO todos (id, owner, status, title, content, due_datetime) VALUES (?, ?, ?, ?, ?, ?)`,
		t.ID, t.Owner, t.Status, t.Title, t.Content, t.DueDatetime,
	)
	return err
}

func (r *TodoRepo) GetByIDOwner(ctx context.Context, id, owner string) (Todo, error) {
	var t Todo
	row := r.db.QueryRowContext(ctx,
		`SELECT id, owner, status, title, content, due_datetime, created_at, updated_at FROM todos WHERE id = ? AND owner = ?`,
		id, owner,
	)
	if err := row.Scan(&t.ID, &t.Owner, &t.Status, &t.Title, &t.Content, &t.DueDatetime, &t.CreatedAt, &t.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Todo{}, sql.ErrNoRows
		}
		return Todo{}, err
	}
	return t, nil
}

func (r *TodoRepo) UpdateByIDOwner(ctx context.Context, id, owner string, title, content, status *string, dueDatetime *time.Time) error {
	t, err := r.GetByIDOwner(ctx, id, owner)
	if err != nil {
		return err
	}
	if title != nil {
		t.Title = *title
	}
	if content != nil {
		t.Content = *content
	}
	if status != nil {
		t.Status = *status
	}
	if dueDatetime != nil {
		t.DueDatetime = dueDatetime
	}
	_, err = r.db.ExecContext(ctx,
		`UPDATE todos SET status = ?, title = ?, content = ?, due_datetime = ? WHERE id = ? AND owner = ?`,
		t.Status, t.Title, t.Content, t.DueDatetime, id, owner,
	)
	return err
}

func (r *TodoRepo) DeleteByIDOwner(ctx context.Context, id, owner string) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM todos WHERE id = ? AND owner = ?`, id, owner)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *TodoRepo) ListByOwner(ctx context.Context, owner string) ([]Todo, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, owner, status, title, content, due_datetime, created_at, updated_at
		 FROM todos WHERE owner = ? ORDER BY created_at DESC`,
		owner,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Todo
	for rows.Next() {
		var t Todo
		if err := rows.Scan(&t.ID, &t.Owner, &t.Status, &t.Title, &t.Content, &t.DueDatetime, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}


