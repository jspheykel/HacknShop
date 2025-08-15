package handlers

import (
	"context"
	"database/sql"
)

type ReportHandler struct{ DB *sql.DB }

type UserSpend struct {
	Username   string
	SpendCents int
}
type DayRevenue struct {
	Day          string
	OrdersCount  int
	RevenueCents int
}
type LowStockItem struct {
	GameID int64
	Title  string
	Stock  int
}

func NewReportHandler(db *sql.DB) *ReportHandler {
	return &ReportHandler{DB: db}
}

func (r *ReportHandler) TopUsersBySpend(ctx context.Context) ([]UserSpend, error) {
	rows, err := r.DB.QueryContext(ctx, `
        SELECT u.username, SUM(o.total_cents) AS spend
        FROM orders o
        JOIN users u ON u.id = o.user_id
        WHERE o.status='PAID'
        GROUP BY u.id, u.username
        ORDER BY spend DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []UserSpend
	for rows.Next() {
		var u UserSpend
		if err := rows.Scan(&u.Username, &u.SpendCents); err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, nil
}

func (r *ReportHandler) RevenuePerDay(ctx context.Context) ([]DayRevenue, error) {
	rows, err := r.DB.QueryContext(ctx, `
        SELECT DATE_FORMAT(o.created_at, '%Y-%m-%d') AS day,
               COUNT(*) AS orders_count,
               COALESCE(SUM(o.total_cents), 0) AS revenue
        FROM orders o
        WHERE o.status = 'PAID'
        GROUP BY day
        ORDER BY day DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []DayRevenue
	for rows.Next() {
		var d DayRevenue
		if err := rows.Scan(&d.Day, &d.OrdersCount, &d.RevenueCents); err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, nil
}

func (r *ReportHandler) LowStock(ctx context.Context, threshold int) ([]LowStockItem, error) {
	rows, err := r.DB.QueryContext(ctx, `
        SELECT id, title, stock FROM games
        WHERE stock <= ? AND is_active=TRUE
        ORDER BY stock ASC, title`, threshold)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []LowStockItem
	for rows.Next() {
		var g LowStockItem
		if err := rows.Scan(&g.GameID, &g.Title, &g.Stock); err != nil {
			return nil, err
		}
		out = append(out, g)
	}
	return out, nil
}
