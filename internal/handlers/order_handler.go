package handlers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type OrderHandler struct {
	DB *sql.DB
}

func NewOrderHandler(db *sql.DB) *OrderHandler {
	return &OrderHandler{DB: db}
}

func (r *OrderHandler) Checkout(ctx context.Context, userID, cartID int64) (orderID int64, totalCents int, err error) {
	tx, err := r.DB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return 0, 0, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	var status string
	if err = tx.QueryRowContext(ctx, `SELECT status FROM carts WHERE id=? AND user_id=?`, cartID, userID).Scan(&status); err != nil {
		if err == sql.ErrNoRows {
			return 0, 0, errors.New("cart not found")
		}
		return 0, 0, err
	}
	if status != "OPEN" {
		return 0, 0, errors.New("cart is not OPEN")
	}

	rows, err := tx.QueryContext(ctx, `
	    SELECT ci.game_id, g.title, ci.qty, g.price_cents, (ci.qty * g.price_cents) AS subtotal, g.stock
        FROM cart_items ci
        JOIN games g ON g.id = ci.game_id
        WHERE ci.cart_id = ?
        FOR UPDATE`, cartID)
	if err != nil {
		return 0, 0, err
	}
	defer rows.Close()

	type rowItem struct {
		gameID   int64
		title    string
		qty      int
		price    int
		subtotal int
		stock    int
	}
	var items []rowItem
	for rows.Next() {
		var it rowItem
		if err = rows.Scan(&it.gameID, &it.title, &it.qty, &it.price, &it.subtotal, &it.stock); err != nil {
			return 0, 0, err
		}
		items = append(items, it)
		totalCents += it.subtotal
	}
	if len(items) == 0 {
		return 0, 0, errors.New("cart is empty")
	}

	// Validate stock
	for _, it := range items {
		if it.qty <= 0 {
			return 0, 0, fmt.Errorf("invalid qty for %s", it.title)
		}
		if it.stock < it.qty {
			return 0, 0, fmt.Errorf("insufficient stock for %s (have %d, need %d)", it.title, it.stock, it.qty)
		}
	}

	// Create order
	res, err := tx.ExecContext(ctx, `
        INSERT INTO orders (user_id, total_cents, status) VALUES (?, ?, 'PAID')`, userID, totalCents)
	if err != nil {
		return 0, 0, err
	}
	orderID, err = res.LastInsertId()
	if err != nil {
		return 0, 0, err
	}

	// Insert order_items + decrement stock
	oiStmt, err := tx.PrepareContext(ctx, `
        INSERT INTO order_items (order_id, game_id, qty, price_cents, subtotal_cents)
        VALUES (?,?,?,?,?)`)
	if err != nil {
		return 0, 0, err
	}
	defer oiStmt.Close()

	decStmt, err := tx.PrepareContext(ctx, `
        UPDATE games SET stock = stock - ? WHERE id = ?`)
	if err != nil {
		return 0, 0, err
	}
	defer decStmt.Close()

	for _, it := range items {
		if _, err = oiStmt.ExecContext(ctx, orderID, it.gameID, it.qty, it.price, it.subtotal); err != nil {
			return 0, 0, err
		}
		// Stock is already locked by FOR UPDATE; subtract now
		res, err := decStmt.ExecContext(ctx, it.qty, it.gameID)
		if err != nil {
			return 0, 0, err
		}
		affected, _ := res.RowsAffected()
		if affected == 0 {
			return 0, 0, fmt.Errorf("failed to update stock for game %d", it.gameID)
		}
	}

	// Clear cart items & close cart
	if _, err = tx.ExecContext(ctx, `DELETE FROM cart_items WHERE cart_id = ?`, cartID); err != nil {
		return 0, 0, err
	}
	if _, err = tx.ExecContext(ctx, `UPDATE carts SET status='CHECKED_OUT' WHERE id=?`, cartID); err != nil {
		return 0, 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, 0, err
	}
	return orderID, totalCents, nil

}
