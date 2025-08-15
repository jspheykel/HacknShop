package models

type Game struct {
	ID         int64
	Title      string
	PriceCents int
	Stock      int
}

type GameAdminView struct {
	ID         int64
	Title      string
	Category   string
	CategoryID int64
	PriceCents int
	Stock      int
	IsActive   bool
}
