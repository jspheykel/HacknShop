package main

import (
	"context"
	"errors"
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
		fmt.Println(util.Red+"Error:"+util.Reset, err)
		return
	}
	fmt.Println(util.Blue + util.Bold + "=== Categories ===" + util.Reset)
	for _, c := range cats {
		fmt.Printf(util.Blue+"[%d] %s - %s\n"+util.Reset, c.ID, c.Name, c.Description)
	}
}

// List Games By Category Controller //
func (u *userActions) ListGamesByCategory() {
	id, err := util.PromptInt("Enter category ID: ")
	if err != nil {
		fmt.Println(util.Red + "Invalid input" + util.Reset)
		return
	}
	games, err := u.GameHandler.ListByCategory(context.Background(), int64(id))
	if err != nil {
		fmt.Println(util.Red+"Error:"+util.Reset, err)
		return
	}
	if len(games) == 0 {
		fmt.Println("No games found")
		return
	}
	fmt.Println(util.Blue + util.Bold + "=== Games ===" + util.Reset)
	for _, g := range games {
		fmt.Printf(util.Blue+"[%d] %s - $%.2f (stock %d)\n"+util.Reset, g.ID, g.Title, float64(g.PriceCents)/100, g.Stock)
	}
}

// Search Games Controller //
func (u *userActions) SearchGames() {
	term := util.Prompt("Search term: ")
	games, err := u.GameHandler.Search(context.Background(), term)
	if err != nil {
		fmt.Println(util.Red+"Error:"+util.Reset, err)
		return
	}
	if len(games) == 0 {
		fmt.Println("No matches")
		return
	}
	fmt.Println(util.Blue + util.Bold + "=== Search Results ===" + util.Reset)
	for _, g := range games {
		fmt.Printf(util.Blue+"[%d] %s - $%.2f (stock %d)\n"+util.Reset, g.ID, g.Title, float64(g.PriceCents)/100, g.Stock)
	}
}

// Add To Cart Controller //
func (u *userActions) AddToCart() {
	gameID, err := util.PromptInt("Game ID: ")
	if err != nil {
		fmt.Println(util.Red + "Invalid input" + util.Reset)
		return
	}
	qty, err := util.PromptInt("Qty: ")
	if err != nil || qty <= 0 {
		fmt.Println(util.Red + "Invalid qty" + util.Reset)
		return
	}

	ctx := context.Background()
	cart, err := u.CartHandler.GetOpenCart(ctx, u.UserID)
	if err != nil {
		fmt.Println(util.Red+"Error:"+util.Reset, err)
		return
	}
	if cart == nil {
		cartID, err := u.CartHandler.CreateCart(ctx, u.UserID)
		if err != nil {
			fmt.Println(util.Red+"Error:"+util.Reset, err)
			return
		}
		cart = &models.Cart{ID: cartID}
	}
	if err := u.CartHandler.AddOrUpdateItem(ctx, cart.ID, int64(gameID), qty); err != nil {
		fmt.Println(util.Red+"Error:"+util.Reset, err)
		return
	}
	fmt.Println(util.Green + "Added to cart." + util.Reset)

}

// View Cart Controller //
func (u *userActions) ViewCart() {
	ctx := context.Background()
	cart, err := u.CartHandler.GetOpenCart(ctx, u.UserID)
	if err != nil {
		fmt.Println(util.Red+"Error:"+util.Reset, err)
		return
	}
	if cart == nil {
		fmt.Println("Cart is empty.")
		return
	}
	items, err := u.CartHandler.ListItems(ctx, cart.ID)
	if err != nil {
		fmt.Println(util.Red+"Error:"+util.Reset, err)
		return
	}
	if len(items) == 0 {
		fmt.Println("Cart is empty.")
		return
	}
	total := 0
	fmt.Println(util.Bold + util.Blue + "=== 🛒 Your Cart ===" + util.Reset)
	for _, v := range items {
		fmt.Printf(util.Blue+"%s x%d = $%.2f\n"+util.Reset, v.Title, v.Qty, float64(v.SubtotalCents)/100)
		total += v.SubtotalCents
	}
	fmt.Printf(util.Cyan+util.Bold+"Total: $%.2f\n"+util.Reset, float64(total)/100)
}

// Checkout Controller //
func (u *userActions) Checkout() {
	ctx := context.Background()
	cart, err := u.CartHandler.GetOpenCart(ctx, u.UserID)
	if err != nil {
		fmt.Println(util.Red+"Error:"+util.Reset, err)
		return
	}
	if cart == nil {
		fmt.Println("No open cart.")
		return
	}

	items, err := u.CartHandler.ListItems(ctx, cart.ID)
	if err != nil {
		fmt.Println(util.Red+"Error:"+util.Reset, err)
		return
	}
	if len(items) == 0 {
		fmt.Println(util.Bold + "🛒 Cart empty. Type [4] to by more games 🎮!" + util.Reset)
		return
	}

	total := 0
	fmt.Println(util.Blue + util.Bold + "=== Checkout Review ===" + util.Reset)
	for _, v := range items {
		fmt.Printf(util.Blue+"%s x%d = $%.2f\n"+util.Reset, v.Title, v.Qty, float64(v.SubtotalCents)/100)
		total += v.SubtotalCents
	}
	fmt.Printf(util.Cyan+util.Bold+"Total: $%.2f\n"+util.Reset, float64(total)/100)
	confirm := util.Prompt(util.Bold + "Proceed? (y/n): " + util.Reset)
	if confirm != "y" && confirm != "Y" {
		fmt.Println("Cancelled.")
		return
	}

	orderID, paidTotal, err := u.OrderHandler.Checkout(ctx, u.UserID, cart.ID)
	if err != nil {
		fmt.Println(util.Red+"Checkout failed:"+util.Reset, err)
		return
	}
	fmt.Printf(util.Green+"Success! Order #%d placed. Total $%.2f\n"+util.Reset, orderID, float64(paidTotal)/100)
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
		fmt.Println(util.Red+"Error:"+util.Reset, err)
		return
	}
	if len(games) == 0 {
		fmt.Println(util.Red + "No games found." + util.Reset)
		return
	}

	// table header
	fmt.Printf(util.Magenta+util.Bold+"\n%-4s  %-30s  %-12s  %-14s  %-10s  %-7s  %-7s\n"+util.Reset,
		"ID", "Title", "Cat ID", "Category", "Price", "Stock", "Active")
	fmt.Println(strings.Repeat("-", 4+2+30+2+12+2+14+2+10+2+7+2+7))

	// rows
	for _, g := range games {
		price := fmt.Sprintf(util.Magenta+"$%.2f"+util.Reset, float64(g.PriceCents)/100)
		active := "No"
		if g.IsActive {
			active = "Yes"
		}
		fmt.Printf(util.Magenta+"%-4d  %-30.30s  %-12d  %-14.14s  %-10s  %-7d  %-7s\n"+util.Reset,
			g.ID, g.Title, g.CategoryID, g.Category, price, g.Stock, active)
	}
}

// Add Game Controller //
func (a *adminActions) AddGame() {
	title := util.Prompt("Title: ")
	catID, err := util.PromptInt("Category ID: ")
	if err != nil {
		fmt.Println(util.Red + "Invalid category" + util.Reset)
		return
	}
	desc := util.Prompt("Description: ")
	price, err := util.PromptInt("Price (in cents): ")
	if err != nil {
		fmt.Println(util.Red + "Invalid price" + util.Reset)
		return
	}
	stock, err := util.PromptInt("Stock: ")
	if err != nil {
		fmt.Println(util.Red + "Invalid stock" + util.Reset)
		return
	}

	id, err := a.GameHandler.AddGame(context.Background(), title, int64(catID), desc, price, stock)
	if err != nil {
		fmt.Println(util.Red+"Error:"+util.Reset, err)
		return
	}
	fmt.Println(util.Green+"Game added with ID:"+util.Reset, id)
}

// Update Stock Price Controller //
func (a *adminActions) UpdateStockPrice() {
	gameID, err := util.PromptInt("Game ID: ")
	if err != nil {
		fmt.Println(util.Red + "Invalid ID" + util.Reset)
		return
	}
	stock, err := util.PromptInt("New Stock: ")
	if err != nil {
		fmt.Println(util.Red + "Invalid stock" + util.Reset)
		return
	}
	price, err := util.PromptInt("New Price (in cents): ")
	if err != nil {
		fmt.Println(util.Red + "Invalid stock" + util.Reset)
		return
	}

	if err := a.GameHandler.UpdateStockPrice(context.Background(), gameID, stock, price); err != nil {
		fmt.Println(util.Red+"Error:"+util.Reset, err)
		return
	}
	fmt.Println(util.Green + "Game updated." + util.Reset)
}

// Delete Game Controller //
func (a *adminActions) DeleteGame() {
	gameID, err := util.PromptInt("Game ID: ")
	if err != nil {
		fmt.Println(util.Red + "Invalid ID" + util.Reset)
		return
	}
	mode := util.Prompt("Delete mode: (soft/hard): ")
	hard := mode == "hard"

	if err := a.GameHandler.DeleteGame(context.Background(), int64(gameID), hard); err != nil {
		fmt.Println(util.Red+"Error:"+util.Reset, err)
		return
	}
	fmt.Println(util.Green + "Game deleted." + util.Reset)
}

// User Reports Controller //
func (a *adminActions) UserReports() {
	res, err := a.ReportHandler.TopUsersBySpend(context.Background())
	if err != nil {
		fmt.Println(util.Red+"Error:"+util.Reset, err)
		return
	}
	fmt.Println(util.Magenta + util.Bold + "=== Top Users by Spend ===" + util.Reset)
	for _, u := range res {
		fmt.Printf(util.Magenta+"%s - $%.2f\n"+util.Reset, u.Username, float64(u.SpendCents)/100)
	}
}

// Order Reports Controller //
func (a *adminActions) OrderReports() {
	res, err := a.ReportHandler.RevenuePerDay(context.Background())
	if err != nil {
		fmt.Println(util.Red+"Error:"+util.Reset, err)
		return
	}
	if len(res) == 0 {
		fmt.Println("No paid orders yet.")
		return
	}

	// Header
	fmt.Printf(util.Magenta+util.Bold+"\n%-12s  %-6s  %-12s\n"+util.Reset, "Date", "Orders", "Revenue")
	fmt.Println(strings.Repeat("-", 12+3+6+3+12))

	var totalOrders int
	var totalRevenue int
	for _, d := range res {
		fmt.Printf(util.Magenta+"%-12s  %-6d  $%-.2f\n"+util.Reset, d.Day, d.OrdersCount, float64(d.RevenueCents)/100)
		totalOrders += d.OrdersCount
		totalRevenue += d.RevenueCents
	}

	fmt.Println(strings.Repeat("-", 12+3+6+3+12))
	fmt.Printf(util.Magenta+"%-12s  %-6d  $%-.2f\n"+util.Reset, "TOTAL", totalOrders, float64(totalRevenue)/100)
}

// Stock Reports Controller //
func (a *adminActions) StockReports() {
	threshold, err := util.PromptInt("Stock threshold: ")
	if err != nil {
		fmt.Println(util.Red + "Invalid number" + util.Reset)
		return
	}
	res, err := a.ReportHandler.LowStock(context.Background(), threshold)
	if err != nil {
		fmt.Println(util.Red+"Error:"+util.Reset, err)
		return
	}
	fmt.Println(util.Magenta + util.Bold + "=== Low Stock Games ===" + util.Reset)
	for _, g := range res {
		fmt.Printf(util.Magenta+"[%d] %s - Stock: %d\n"+util.Reset, g.GameID, g.Title, g.Stock)
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
	// fmt.Println("Note: ensure your DB has bcrypt password hashes for login.")

	for {
		session, err := cli.LoginOrRegister(auth)
		if err != nil {
			if errors.Is(err, cli.ErrAppExit) {
				fmt.Println(util.Green + "Thankyou for trusting HacknShop Games Store 👾. See you later 👋🏻!" + util.Reset)
				return
			}
			fmt.Println(util.Red+"Error:"+util.Reset, err)
			continue
		}

		fmt.Printf(util.Blue+util.Bold+"Hello, %s! Welcome to HacknShop Games Store 👾!\n"+util.Reset, session.Name)

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

	}

}
