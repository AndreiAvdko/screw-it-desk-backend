package main

import (
	"context"
	"flag"
	"log"
	"screw-it-desk-backend/services/users/internal/app"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

// Спагетти-код с хэндлерами

type UserResponse struct {
	ID        string  `json:"id"`
	Email     string  `json:"email"`
	Username  string  `json:"username"`
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
	Phone     *string `json:"phone"`
	Role      string  `json:"role"`
	Status    string  `json:"status"`
}

type CustomerResponse struct {
	ID           string  `json:"id"`
	Email        string  `json:"email"`
	Username     string  `json:"username"`
	FirstName    *string `json:"firstName"`
	LastName     *string `json:"lastName"`
	Phone        *string `json:"phone"`
	Role         string  `json:"role"`
	Status       string  `json:"status"`
	Organization string  `json:"organization"`
	TicketsCount int     `json:"ticketsCount"`
}

type UsersListResponse struct {
	Data []UserResponse `json:"data"`
}

type CustomersListResponse struct {
	Data       []CustomerResponse `json:"data"`
	Pagination Pagination         `json:"pagination"`
}

type Pagination struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
	Pages int `json:"pages"`
}

// Чистый main с инициализацией app
func main() {
	ctx := context.Background()

	a, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("failed to init app: %s", err.Error())
	}

	err = a.Run()
	if err != nil {
		log.Fatalf("failed to run app: %s", err.Error())
	}
}
