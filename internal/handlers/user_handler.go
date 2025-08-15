package handlers

import (
	"context"
	"database/sql"

	"github.com/jspheykel/HacknShop/internal/models"
)

type UserHandler struct{ DB *sql.DB }

func NewUserHandler(db *sql.DB) *UserHandler {
	return &UserHandler{DB: db}
}

func (r *UserHandler) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	row := r.DB.QueryRowContext(ctx, `
        SELECT id, username, email, password_hash, is_admin, created_at, updated_at
        FROM users WHERE username = ?`, username)
	u := models.User{}
	if err := row.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.IsAdmin, &u.CreatedAt, &u.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserHandler) Create(ctx context.Context, username, email, pwdHash string, isAdmin bool) (int64, error) {
	res, err := r.DB.ExecContext(ctx, `
        INSERT INTO users (username, email, password_hash, is_admin)
        VALUES (?,?,?,?)`, username, email, pwdHash, isAdmin)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}
