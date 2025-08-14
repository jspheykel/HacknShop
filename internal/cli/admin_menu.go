package cli

import (
	"fmt"

	"github.com/jspheykel/HacknShop/internal/util"
)

type AdminActions interface {
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
		fmt.Println("1) Add Games")
		fmt.Println("2) Update Stocks & Price")
		fmt.Println("3) Delete Games")
		fmt.Println("4) User Reports")
		fmt.Println("5) Order Reports")
		fmt.Println("6) Stock Reports")
		fmt.Println("7) Exit")
		choice := util.Prompt("Choose: ")

		switch choice {
		case "1":
			act.AddGame()
		case "2":
			act.UpdateStockPrice()
		case "3":
			act.DeleteGame()
		case "4":
			act.UserReports()
		case "5":
			act.OrderReports()
		case "6":
			act.StockReports()
		case "7":
			return
		default:
			fmt.Println("Invalid choice.")
		}
	}
}
