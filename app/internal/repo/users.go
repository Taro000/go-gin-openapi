package repo

import (
	"context"
	"database/sql"
	"errors"
)

type User struct {
	UID      string
	Nickname string
	Email    string
}

type UserRepo struct {
	db *sql.DB
}

func (r *UserRepo) Create(ctx context.Context, u User) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO users (uid, nickname, email) VALUES (?, ?, ?)`,
		u.UID, u.Nickname, u.Email,
	)
	return err
}

func (r *UserRepo) GetByUID(ctx context.Context, uid string) (User, error) {
	var u User
	row := r.db.QueryRowContext(ctx, `SELECT uid, nickname, email FROM users WHERE uid = ?`, uid)
	if err := row.Scan(&u.UID, &u.Nickname, &u.Email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, sql.ErrNoRows
		}
		return User{}, err
	}
	return u, nil
}

func (r *UserRepo) Update(ctx context.Context, uid string, nickname *string, email *string) (User, error) {
	// Fetch existing and apply partial updates
	u, err := r.GetByUID(ctx, uid)
	if err != nil {
		return User{}, err
	}
	if nickname != nil {
		u.Nickname = *nickname
	}
	if email != nil {
		u.Email = *email
	}
	_, err = r.db.ExecContext(ctx, `UPDATE users SET nickname = ?, email = ? WHERE uid = ?`, u.Nickname, u.Email, uid)
	if err != nil {
		return User{}, err
	}
	return u, nil
}


