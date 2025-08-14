package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/jspheykel/HacknShop/internal/cli"
	"github.com/jspheykel/HacknShop/internal/config"
	"github.com/jspheykel/HacknShop/internal/db"
	"github.com/jspheykel/HacknShop/internal/handlers"
	"github.com/jspheykel/HacknShop/internal/models"
	"github.com/jspheykel/HacknShop/internal/service"
	"github.com/jspheykel/HacknShop/internal/util"
)

// User Menu //
type userActions struct {
	GameHandler  *handlers.GameHandler
	CartHandler  *handlers.CartHandler
	OrderHandler *handlers.OrderHandler
	UserID       int64
}

// List Categories Controller //
func (u *userActions) ListCategories() {
	cats, err := u.GameHandler.ListCategories(context.Background())
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("=== Categories ===")
	for _, c := range cats {
		fmt.Printf("[%d] %s - %s\n", c.ID, c.Name, c.Description)
	}
}

// List Games By Category Controller //
func (u *userActions) ListGamesByCategory() {
	id, err := util.PromptInt("Enter category ID: ")
	if err != nil {
		fmt.Println("Invalid input")
		return
	}
	games, err := u.GameHandler.ListByCategory(context.Background(), int64(id))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if len(games) == 0 {
		fmt.Println("No games found")
		return
	}
	fmt.Println("=== Games ===")
	for _, g := range games {
		fmt.Printf("[%d] %s - $%.2f (stock %d)\n", g.ID, g.Title, float64(g.PriceCents)/100, g.Stock)
	}
}

// Search Games Controller //
func (u *userActions) SearchGames() {
	term := util.Prompt("Search term: ")
	games, err := u.GameHandler.Search(context.Background(), term)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if len(games) == 0 {
		fmt.Println("No matches")
		return
	}
	fmt.Println("=== Search Results ===")
	for _, g := range games {
		fmt.Printf("[%d] %s - $%.2f (stock %d)\n", g.ID, g.Title, float64(g.PriceCents)/100, g.Stock)
	}
}

// Add To Cart Controller //
func (u *userActions) AddToCart() {
	gameID, err := util.PromptInt("Game ID: ")
	if err != nil {
		fmt.Println("Invalid input")
		return
	}
	qty, err := util.PromptInt("Qty: ")
	if err != nil || qty <= 0 {
		fmt.Println("Invalid qty")
		return
	}

	ctx := context.Background()
	cart, err := u.CartHandler.GetOpenCart(ctx, u.UserID)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if cart == nil {
		cartID, err := u.CartHandler.CreateCart(ctx, u.UserID)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		cart = &models.Cart{ID: cartID}
	}
	if err := u.CartHandler.AddOrUpdateItem(ctx, cart.ID, int64(gameID), qty); err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Added to cart.")

}

// View Cart Controller //
func (u *userActions) ViewCart() {
	ctx := context.Background()
	cart, err := u.CartHandler.GetOpenCart(ctx, u.UserID)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if cart == nil {
		fmt.Println("Cart is empty.")
		return
	}
	items, err := u.CartHandler.ListItems(ctx, cart.ID)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if len(items) == 0 {
		fmt.Println("Cart is empty.")
		return
	}
	total := 0
	fmt.Println("=== Your Cart ===")
	for _, v := range items {
		fmt.Printf("%s x%d = $%.2f\n", v.Title, v.Qty, float64(v.SubtotalCents)/100)
		total += v.SubtotalCents
	}
	fmt.Printf("Total: $%.2f\n", float64(total)/100)
}

// Checkout Controller //
func (u *userActions) Checkout() {
	ctx := context.Background()
	cart, err := u.CartHandler.GetOpenCart(ctx, u.UserID)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if cart == nil {
		fmt.Println("No open cart.")
		return
	}

	items, err := u.CartHandler.ListItems(ctx, cart.ID)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if len(items) == 0 {
		fmt.Println("Cart empty.")
		return
	}

	total := 0
	fmt.Println("=== Checkout Review ===")
	for _, v := range items {
		fmt.Printf("%s x%d = $%.2f\n", v.Title, v.Qty, float64(v.SubtotalCents)/100)
		total += v.SubtotalCents
	}
	fmt.Printf("Total: $%.2f\n", float64(total)/100)
	confirm := util.Prompt("Proceed? (y/n): ")
	if confirm != "y" && confirm != "Y" {
		fmt.Println("Cancelled.")
		return
	}

	orderID, paidTotal, err := u.OrderHandler.Checkout(ctx, u.UserID, cart.ID)
	if err != nil {
		fmt.Println("Checkout failed:", err)
		return
	}
	fmt.Printf("Success! Order #%d placed. Total $%.2f\n", orderID, float64(paidTotal)/100)
}

// Admin Menu //
type adminActions struct {
	GameHandler   *handlers.GameHandler
	ReportHandler *handlers.ReportHandler
}

// Show All List Of Game
func (a *adminActions) ListAllGames() {
	ctx := context.Background()
	games, err := a.GameHandler.ListAllGames(ctx)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if len(games) == 0 {
		fmt.Println("No games found.")
		return
	}

	// table header
	fmt.Printf("\n%-4s  %-30s  %-14s  %-10s  %-7s  %-7s\n", "ID", "Title", "Category", "Price", "Stock", "Active")
	fmt.Println(strings.Repeat("-", 4+2+30+2+14+2+10+2+7+2+7))

	// rows
	for _, g := range games {
		price := fmt.Sprintf("$%.2f", float64(g.PriceCents)/100)
		active := "No"
		if g.IsActive {
			active = "Yes"
		}
		fmt.Printf("%-4d  %-30.30s  %-14.14s  %-10s  %-7d  %-7s\n",
			g.ID, g.Title, g.Category, price, g.Stock, active)
	}
}

// Add Game Controller //
func (a *adminActions) AddGame() {
	title := util.Prompt("Title: ")
	catID, err := util.PromptInt("Category ID: ")
	if err != nil {
		fmt.Println("Invalid category")
		return
	}
	desc := util.Prompt("Description: ")
	price, err := util.PromptInt("Price (in cents): ")
	if err != nil {
		fmt.Println("Invalid price")
		return
	}
	stock, err := util.PromptInt("Stock: ")
	if err != nil {
		fmt.Println("Invalid stock")
		return
	}

	id, err := a.GameHandler.AddGame(context.Background(), title, int64(catID), desc, price, stock)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Game added with ID:", id)
}

// Update Stock Price Controller //
func (a *adminActions) UpdateStockPrice() {
	gameID, err := util.PromptInt("Game ID: ")
	if err != nil {
		fmt.Println("Invalid ID")
		return
	}
	stock, err := util.PromptInt("New Stock: ")
	if err != nil {
		fmt.Println("Invalid stock")
		return
	}
	price, err := util.PromptInt("New Price (in cents): ")
	if err != nil {
		fmt.Println("Invalid price")
		return
	}

	if err := a.GameHandler.UpdateStockPrice(context.Background(), gameID, stock, price); err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Game updated.")
}

// Delete Game Controller //
func (a *adminActions) DeleteGame() {
	gameID, err := util.PromptInt("Game ID: ")
	if err != nil {
		fmt.Println("Invalid ID")
		return
	}
	mode := util.Prompt("Delete mode: (soft/hard): ")
	hard := mode == "hard"

	if err := a.GameHandler.DeleteGame(context.Background(), int64(gameID), hard); err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Game deleted.")
}

// User Reports Controller //
func (a *adminActions) UserReports() {
	res, err := a.ReportHandler.TopUsersBySpend(context.Background())
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("=== Top Users by Spend ===")
	for _, u := range res {
		fmt.Printf("%s - $%.2f\n", u.Username, float64(u.SpendCents)/100)
	}
}

// Order Reports Controller //
func (a *adminActions) OrderReports() {
	res, err := a.ReportHandler.RevenuePerDay(context.Background())
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if len(res) == 0 {
		fmt.Println("No paid orders yet.")
		return
	}

	// Header
	fmt.Printf("\n%-12s  %-6s  %-12s\n", "Date", "Orders", "Revenue")
	fmt.Println(strings.Repeat("-", 12+3+6+3+12))

	var totalOrders int
	var totalRevenue int
	for _, d := range res {
		fmt.Printf("%-12s  %-6d  $%-.2f\n", d.Day, d.OrdersCount, float64(d.RevenueCents)/100)
		totalOrders += d.OrdersCount
		totalRevenue += d.RevenueCents
	}

	fmt.Println(strings.Repeat("-", 12+3+6+3+12))
	fmt.Printf("%-12s  %-6d  $%-.2f\n", "TOTAL", totalOrders, float64(totalRevenue)/100)
}

// Stock Reports Controller //
func (a *adminActions) StockReports() {
	threshold, err := util.PromptInt("Stock threshold: ")
	if err != nil {
		fmt.Println("Invalid number")
		return
	}
	res, err := a.ReportHandler.LowStock(context.Background(), threshold)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("=== Low Stock Games ===")
	for _, g := range res {
		fmt.Printf("[%d] %s - Stock: %d\n", g.GameID, g.Title, g.Stock)
	}
}

func main() {
	cfg := config.Default()
	dsn := cfg.DSN()

	sqlDB, err := db.Open(dsn)
	if err != nil {
		panic(err)
	}
	defer sqlDB.Close()

	userHandler := handlers.NewUserHandler(sqlDB)
	auth := &service.AuthService{Users: userHandler}

	// On first run, help seed admin password properly:
	fmt.Println("Note: ensure your DB has bcrypt password hashes for login.")

	session, err := cli.LoginOrRegister(auth)
	if err != nil {
		fmt.Println("Goodbye.")
		return
	}

	fmt.Printf("Hello, %s! (admin=%v)\n", session.Name, session.IsAdmin)

	gameHandler := handlers.NewGameHandler(sqlDB)
	cartHandler := handlers.NewCartHandler(sqlDB)
	orderHandler := handlers.NewOrderHandler(sqlDB)

	if session.IsAdmin {
		acts := &adminActions{
			GameHandler:   gameHandler,
			ReportHandler: handlers.NewReportHandler(sqlDB),
		}
		cli.AdminMenu(acts)
	} else {
		acts := &userActions{
			GameHandler:  gameHandler,
			CartHandler:  cartHandler,
			OrderHandler: orderHandler,
			UserID:       session.UserID,
		}
		cli.UserMenu(acts)
	}

	fmt.Println("Thanks for using Games E‑Commerce CLI!")
}
