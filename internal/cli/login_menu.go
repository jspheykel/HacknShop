package cli

import (
	"context"
	"fmt"

	"github.com/jspheykel/HacknShop/internal/service"
	"github.com/jspheykel/HacknShop/internal/util"
)

type Session struct {
	UserID  int64
	IsAdmin bool
	Name    string
}

func LoginOrRegister(auth *service.AuthService) (*Session, error) {
	for {
		fmt.Println("=== Welcome to HacknShop ===")
		fmt.Println("1) Login")
		fmt.Println("2) Register")
		fmt.Println("3) Exit")
		choice := util.Prompt("Choose: ")

		switch choice {
		case "1":
			username := util.Prompt("Username: ")
			password := util.Prompt("Password: ")
			u, err := auth.Login(context.Background(), username, password)
			if err != nil {
				fmt.Println("Login failed: ", err.Error())
				continue
			}
			return &Session{UserID: u.ID, IsAdmin: u.IsAdmin, Name: u.Username}, nil
		case "2":
			username := util.Prompt("Username: ")
			email := util.Prompt("Email: ")
			password := util.Prompt("Password: ")
			_, err := auth.Register(context.Background(), username, email, password)
			if err != nil {
				fmt.Println("Register failed:", err.Error())
				continue
			}
			fmt.Println("Register success. Please login.")
		case "3":
			return nil, fmt.Errorf("exit")
		default:
			fmt.Println("Invalid choice.")
		}
	}
}
