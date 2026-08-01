package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

// TODO организовать handler как в payment-service репозитории с зависимостью через usecase
var adminsCache []UserResponse

func GetAdminsHandler(w http.ResponseWriter, r *http.Request) {

	if err := loadAdminsFromFile(); err != nil {
		fmt.Printf("Error loading admins.json: %v\n", err)
		adminsCache = []UserResponse{}
	}

	response := UsersListResponse{
		Data: adminsCache,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func loadAdminsFromFile() error {
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("cannot get current directory: %w", err)
	}

	filePath := filepath.Join(currentDir, "/mocks/admins.json")

	byteValue, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("cannot read file: %w", err)
	}

	var rawData struct {
		Admins []UserResponse `json:"admins"`
	}

	if err := json.Unmarshal(byteValue, &rawData); err != nil {
		return fmt.Errorf("json parse error: %w", err)
	}

	adminsCache = rawData.Admins
	fmt.Printf("Loaded %d admins from admins.json\n", len(adminsCache))

	return nil
}

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

type UsersListResponse struct {
	Data []UserResponse `json:"data"`
}
