package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

var customersCache []CustomerResponse

func GetCustomersHandler(w http.ResponseWriter, r *http.Request) {
	if err := loadCustomersFromFile(); err != nil {
		fmt.Printf("Error loading customers.json: %v\n", err)
		customersCache = []CustomerResponse{}
	}
	// Пример получения данных из БД
	// 	// Делаем запрос на получение измененной записи из таблицы note
	// 	builderSelectOne := sq.Select("id", "title", "body", "created_at", "updated_at").
	// 		From("note").
	// 		PlaceholderFormat(sq.Dollar).
	// 		Where(sq.Eq{"id": req.GetId()}).
	// 		Limit(1)

	// 	query, args, err := builderSelectOne.ToSql()
	// 	if err != nil {
	// 		log.Fatalf("failed to build query: %v", err)
	// 	}

	// 	var id int64
	// 	var title, body string
	// 	var createdAt time.Time
	// 	var updatedAt sql.NullTime

	// 	err = s.pool.QueryRow(ctx, query, args...).Scan(&id, &title, &body, &createdAt, &updatedAt)
	// 	if err != nil {
	// 		log.Fatalf("failed to select notes: %v", err)
	// 	}

	// 	log.Printf("id: %d, title: %s, body: %s, created_at: %v, updated_at: %v\n", id, title, body, createdAt, updatedAt)

	// 	var updatedAtTime *timestamppb.Timestamp
	// 	if updatedAt.Valid {
	// 		updatedAtTime = timestamppb.New(updatedAt.Time)
	// 	}

	// 	return &desc.GetResponse{
	// 		Note: &desc.Note{
	// 			Id:        id,
	// 			Info:      &desc.NoteInfo{Title: title, Content: body},
	// 			CreatedAt: timestamppb.New(createdAt),
	// 			UpdatedAt: updatedAtTime,
	// 		},
	// 	}, nil
	// }

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	status := r.URL.Query().Get("status")
	if status == "" {
		status = "all"
	}

	sort := r.URL.Query().Get("sort")
	if sort == "" {
		sort = "createdAt:desc"
	}

	fmt.Printf("GET /v1/users/customers - page: %d, limit: %d, status: %s, sort: %s\n", page, limit, status, sort)

	var filteredCustomers []CustomerResponse

	if status == "all" {
		filteredCustomers = customersCache
	} else {
		for _, customer := range customersCache {
			if customer.Status == status {
				filteredCustomers = append(filteredCustomers, customer)
			}
		}
	}

	total := len(filteredCustomers)
	pages := (total + limit - 1) / limit
	if pages == 0 && total == 0 {
		pages = 1
	}

	if page > pages {
		response := CustomersListResponse{
			Data: []CustomerResponse{},
			Pagination: Pagination{
				Page:  page,
				Limit: limit,
				Total: total,
				Pages: pages,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	start := (page - 1) * limit
	end := start + limit
	if end > total {
		end = total
	}

	paginatedCustomers := filteredCustomers[start:end]

	response := CustomersListResponse{
		Data: paginatedCustomers,
		Pagination: Pagination{
			Page:  page,
			Limit: limit,
			Total: total,
			Pages: pages,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func loadCustomersFromFile() error {
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("cannot get current directory: %w", err)
	}

	filePath := filepath.Join(currentDir, "/customers.json")

	byteValue, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("cannot read file: %w", err)
	}

	var rawData struct {
		Customers []CustomerResponse `json:"customers"`
	}

	if err := json.Unmarshal(byteValue, &rawData); err != nil {
		return fmt.Errorf("json parse error: %w", err)
	}

	customersCache = rawData.Customers
	fmt.Printf("Loaded %d customers from customers.json\n", len(customersCache))

	return nil
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
