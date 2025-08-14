package main

import (
	"fmt"

	"github.com/jspheykel/HacknShop/internal/cli"
	"github.com/jspheykel/HacknShop/internal/config"
	"github.com/jspheykel/HacknShop/internal/db"
	"github.com/jspheykel/HacknShop/internal/handlers"
	"github.com/jspheykel/HacknShop/internal/service"
)

type userActions struct {
}

func (u *userActions) ListCategories()      { fmt.Println("[TODO] list categories from DB") }
func (u *userActions) ListGamesByCategory() { fmt.Println("[TODO] list games by chosen category") }
func (u *userActions) SearchGames()         { fmt.Println("[TODO] search by title LIKE") }
func (u *userActions) AddToCart()           { fmt.Println("[TODO] upsert cart_items for OPEN cart") }
func (u *userActions) ViewCart()            { fmt.Println("[TODO] show cart items & totals") }
func (u *userActions) Checkout()            { fmt.Println("[TODO] tx: create order, move items, dec stock") }

type adminActions struct{}

func (a *adminActions) AddGame()          { fmt.Println("[TODO] insert into games") }
func (a *adminActions) UpdateStockPrice() { fmt.Println("[TODO] update games set stock=.., price=..") }
func (a *adminActions) DeleteGame()       { fmt.Println("[TODO] soft/hard delete game") }
func (a *adminActions) UserReports()      { fmt.Println("[TODO] top users by spend query") }
func (a *adminActions) OrderReports()     { fmt.Println("[TODO] revenue per day query") }
func (a *adminActions) StockReports()     { fmt.Println("[TODO] low stock query") }

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

	if session.IsAdmin {
		acts := &adminActions{}
		cli.AdminMenu(acts)
	} else {
		acts := &userActions{}
		cli.UserMenu(acts)
	}

	fmt.Println("Thanks for using Games E‑Commerce CLI!")
}
