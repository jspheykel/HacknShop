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
		fmt.Println(util.Cyan + util.Bold + "\n=== User Menu ===" + util.Reset)
		fmt.Println(util.Yellow + "1)" + util.Reset + " List Categories")
		fmt.Println(util.Yellow + "2)" + util.Reset + " List Games By Category")
		fmt.Println(util.Yellow + "3)" + util.Reset + " Search Games")
		fmt.Println(util.Yellow + "4)" + util.Reset + " Add to Cart (by Game ID)")
		fmt.Println(util.Yellow + "5)" + util.Reset + " View Cart")
		fmt.Println(util.Yellow + "6)" + util.Reset + " Checkout")
		fmt.Println(util.Yellow + "7)" + util.Reset + " Logout")
		choice := util.Prompt(util.Green + "Choose: " + util.Reset)

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
			fmt.Println(util.Red + "Invalid choice." + util.Reset)
		}
	}
}
