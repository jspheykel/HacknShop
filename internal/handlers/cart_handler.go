package handlers

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jspheykel/HacknShop/internal/models"
)

type CartHandler struct{ DB *sql.DB }

func NewCartHandler(db *sql.DB) *CartHandler { return &CartHandler{DB: db} }

func (r *CartHandler) GetOpenCart(ctx context.Context, userID int64) (*models.Cart, error) {
	row := r.DB.QueryRowContext(ctx, `
        SELECT id, user_id, status FROM carts
        WHERE user_id = ? AND status = 'OPEN' LIMIT 1`, userID)
	var c models.Cart
	if err := row.Scan(&c.ID, &c.UserID, &c.Status); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (r *CartHandler) CreateCart(ctx context.Context, userID int64) (int64, error) {
	res, err := r.DB.ExecContext(ctx, `INSERT INTO carts (user_id, status) VALUES (?, 'OPEN')`, userID)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *CartHandler) AddOrUpdateItem(ctx context.Context, cartID, gameID int64, qty int) error {
	// Upsert pattern
	_, err := r.DB.ExecContext(ctx, `
        INSERT INTO cart_items (cart_id, game_id, qty)
        VALUES (?,?,?)
        ON DUPLICATE KEY UPDATE qty = qty + VALUES(qty)`,
		cartID, gameID, qty)
	return err
}

func (r *CartHandler) ListItems(ctx context.Context, cartID int64) ([]models.CartItemView, error) {
	rows, err := r.DB.QueryContext(ctx, `
        SELECT ci.game_id, g.title, ci.qty, g.price_cents, (ci.qty * g.price_cents) AS subtotal
        FROM cart_items ci
        JOIN games g ON g.id = ci.game_id
        WHERE ci.cart_id = ?`, cartID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.CartItemView
	for rows.Next() {
		var v models.CartItemView
		if err := rows.Scan(&v.GameID, &v.Title, &v.Qty, &v.PriceCents, &v.SubtotalCents); err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, nil
}

func (r *CartHandler) ClearCart(ctx context.Context, cartID int64) error {
	_, err := r.DB.ExecContext(ctx, `DELETE FROM cart_items WHERE cart_id = ?`, cartID)
	return err
}

func (r *CartHandler) CloseCart(ctx context.Context, cartID int64) error {
	res, err := r.DB.ExecContext(ctx, `UPDATE carts SET status='CHECKED_OUT' WHERE id=?`, cartID)
	if err != nil {
		return err
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		return errors.New("no cart updated")
	}
	return nil
}
