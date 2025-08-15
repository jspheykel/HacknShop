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
		fmt.Println(util.Magenta + util.Bold + "\n=== Admin Menu ===" + util.Reset)
		fmt.Println(util.Yellow + "1)" + util.Reset + " List All Games (Table)")
		fmt.Println(util.Yellow + "2)" + util.Reset + " Add Games")
		fmt.Println(util.Yellow + "3)" + util.Reset + " Update Stocks & Price")
		fmt.Println(util.Yellow + "4)" + util.Reset + " Delete Games")
		fmt.Println(util.Yellow + "5)" + util.Reset + " User Reports")
		fmt.Println(util.Yellow + "6)" + util.Reset + " Order Reports")
		fmt.Println(util.Yellow + "7)" + util.Reset + " Stock Reports")
		fmt.Println(util.Yellow + "8)" + util.Reset + " Logout")
		choice := util.Prompt(util.Green + "Choose: " + util.Reset)

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
			fmt.Println(util.Red + "Invalid choice." + util.Reset)
		}
	}
}
