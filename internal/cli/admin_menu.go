package cli

import (
	"fmt"

	"github.com/jspheykel/HacknShop/internal/util"
)

type AdminActions interface {
	ListAllGames()
	AddGame()
	UpdateStockPrice()
	DeleteGame()
	UserReports()
	OrderReports()
	StockReports()
}

func AdminMenu(act AdminActions) {
	for {
		fmt.Println("\n=== Admin Menu ===")
		fmt.Println("1) List All Games (Table)")
		fmt.Println("2) Add Games")
		fmt.Println("3) Update Stocks & Price")
		fmt.Println("4) Delete Games")
		fmt.Println("5) User Reports")
		fmt.Println("6) Order Reports")
		fmt.Println("7) Stock Reports")
		fmt.Println("8) Exit")
		choice := util.Prompt("Choose: ")

		switch choice {
		case "1":
			act.ListAllGames()
		case "2":
			act.AddGame()
		case "3":
			act.UpdateStockPrice()
		case "4":
			act.DeleteGame()
		case "5":
			act.UserReports()
		case "6":
			act.OrderReports()
		case "7":
			act.StockReports()
		case "8":
			return
		default:
			fmt.Println("Invalid choice.")
		}
	}
}
