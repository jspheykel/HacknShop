package models

type Cart struct {
	ID     int64
	UserID int64
	Status string
}

type CartItemView struct {
	GameID        int64
	Title         string
	Qty           int
	PriceCents    int
	SubtotalCents int
}
