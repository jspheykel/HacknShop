package cli

import (
	"fmt"

	"github.com/jspheykel/HacknShop/internal/util"
)

type UserActions interface {
	ListCategories()
	ListGamesByCategory()
	SearchGames()
	AddToCart()
	ViewCart()
	Checkout()
}

func UserMenu(act UserActions) {
	for {
		fmt.Println("\n=== User Menu ===")
		fmt.Println("1) List Categories")
		fmt.Println("2) List Games By Category")
		fmt.Println("3) Search Games")
		fmt.Println("4) Add to Cart (by Game ID)")
		fmt.Println("5) View Cart")
		fmt.Println("6) Checkout")
		fmt.Println("7) Exit")
		choice := util.Prompt("Choose: ")

		switch choice {
		case "1":
			act.ListCategories()
		case "2":
			act.ListGamesByCategory()
		case "3":
			act.SearchGames()
		case "4":
			act.AddToCart()
		case "5":
			act.ViewCart()
		case "6":
			act.Checkout()
		case "7":
			return
		default:
			fmt.Println("Invalid choice.")
		}
	}
}
