package handlers

import (
	"context"
	"database/sql"

	"github.com/jspheykel/HacknShop/internal/models"
)

type GameHandler struct{ DB *sql.DB }

func NewGameHandler(db *sql.DB) *GameHandler {
	return &GameHandler{DB: db}
}

func (r *GameHandler) ListCategories(ctx context.Context) ([]models.Category, error) {
	rows, err := r.DB.QueryContext(ctx, `SELECT id, name, description FROM categories ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Description); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, nil
}

func (r *GameHandler) ListByCategory(ctx context.Context, categoryID int64) ([]models.Game, error) {
	rows, err := r.DB.QueryContext(ctx, `
        SELECT id, title, price_cents, stock FROM games
        WHERE category_id = ? AND is_active = TRUE
        ORDER BY title`, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.Game
	for rows.Next() {
		var g models.Game
		if err := rows.Scan(&g.ID, &g.Title, &g.PriceCents, &g.Stock); err != nil {
			return nil, err
		}
		out = append(out, g)
	}
	return out, nil
}

func (r *GameHandler) Search(ctx context.Context, term string) ([]models.Game, error) {
	rows, err := r.DB.QueryContext(ctx, `
        SELECT id, title, price_cents, stock FROM games
        WHERE title LIKE CONCAT('%',?,'%') AND is_active = TRUE
        ORDER BY title`, term)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.Game
	for rows.Next() {
		var g models.Game
		if err := rows.Scan(&g.ID, &g.Title, &g.PriceCents, &g.Stock); err != nil {
			return nil, err
		}
		out = append(out, g)
	}
	return out, nil
}

func (r *GameHandler) AddGame(ctx context.Context, title string, categoryID int64, desc string, priceCents, stock int) (int64, error) {
	res, err := r.DB.ExecContext(ctx, `
	INSERT INTO games (title, category_id, description, price_cents, stock, is_active)
	VALUE (?, ?, ?, ?, ?, TRUE)`,
		title, categoryID, desc, priceCents, stock)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *GameHandler) UpdateStockPrice(ctx context.Context, gameID, stock, priceCents int) error {
	_, err := r.DB.ExecContext(ctx, `
	UPDATE games SET stock=?, price_cents=? WHERE id=?`, stock, priceCents, gameID)
	return err
}

func (r *GameHandler) DeleteGame(ctx context.Context, gameID int64, hard bool) error {
	if hard {
		_, err := r.DB.ExecContext(ctx, `DELETE FROM games WHERE id=?`, gameID)
		return err
	}
	_, err := r.DB.ExecContext(ctx, `UPDATE games SET is_active=FALSE WHERE id=?`, gameID)
	return err
}

func (r *GameHandler) ListAllGames(ctx context.Context) ([]models.GameAdminView, error) {
	rows, err := r.DB.QueryContext(ctx, `
        SELECT g.id, g.title, c.id AS category_id, c.name AS category, 
               g.price_cents, g.stock, g.is_active
        FROM games g
        JOIN categories c ON c.id = g.category_id
        ORDER BY g.title`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.GameAdminView
	for rows.Next() {
		var v models.GameAdminView
		if err := rows.Scan(&v.ID, &v.Title, &v.CategoryID, &v.Category, &v.PriceCents, &v.Stock, &v.IsActive); err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, nil
}
