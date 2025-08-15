package cli

import (
	"context"
	"errors"
	"fmt"

	"github.com/jspheykel/HacknShop/internal/service"
	"github.com/jspheykel/HacknShop/internal/util"
)

type Session struct {
	UserID  int64
	IsAdmin bool
	Name    string
}

var ErrAppExit = errors.New("HacknShop")

func LoginOrRegister(auth *service.AuthService) (*Session, error) {
	for {
		fmt.Println(util.Magenta + "\n====================================" + util.Reset)
		fmt.Println(util.Magenta + util.Bold + "\n===== 👾 Welcome to HacknShop! =====" + util.Reset)
		fmt.Println(util.Magenta + "\n=== your trusted games store 🎮🕹️ ===" + util.Reset)
		fmt.Println(util.Magenta + "\n====================================" + util.Reset)
		fmt.Println(util.Green + "1) Login" + util.Reset)
		fmt.Println(util.Yellow + "2) Register" + util.Reset)
		fmt.Println(util.Red + "3) Exit" + util.Reset)
		choice := util.Prompt(util.Green + "Choose: " + util.Reset)

		switch choice {
		case "1":
			username := util.Prompt("Username: ")
			password := util.Prompt("Password: ")
			u, err := auth.Login(context.Background(), username, password)
			if err != nil {
				fmt.Println(util.Red+"Login failed: "+util.Reset, err.Error())
				continue
			}
			return &Session{UserID: u.ID, IsAdmin: u.IsAdmin, Name: u.Username}, nil
		case "2":
			username := util.Prompt("Username: ")
			email := util.Prompt("Email: ")
			password := util.Prompt("Password: ")
			_, err := auth.Register(context.Background(), username, email, password)
			if err != nil {
				fmt.Println(util.Red+"Register failed:"+util.Reset, err.Error())
				continue
			}
			fmt.Println(util.Green + "Register success. Please login." + util.Reset)
		case "3":
			return nil, ErrAppExit
		default:
			fmt.Println("Invalid choice.")
		}
	}
}
