package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"screw-it-desk-backend/services/users/internal/config"
	"screw-it-desk-backend/services/users/internal/config/env"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v4/pgxpool"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

// type server struct {
// 	desc.UnimplementedNoteV1Server
// 	pool *pgxpool.Pool
// }

func main() {
	flag.Parse()
	ctx := context.Background()

	// Считываем переменные окружения
	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	httpConfig, err := env.NewHTTPConfig()
	if err != nil {
		log.Fatalf("failed to get grpc config: %v", err)
	}

	pgConfig, err := env.NewPGConfig()
	if err != nil {
		log.Fatalf("failed to get pg config: %v", err)
	}

	// TODO заменить прослушивание порта на запуск http-сервера (используем )
	lis, err := net.Listen("tcp", httpConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Создаем пул соединений с базой данных
	pool, err := pgxpool.Connect(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Пример для запуска grpc-сервера
	// s := grpc.NewServer()
	// reflection.Register(s)
	// desc.RegisterNoteV1Server(s, &server{pool: pool})

	// Запуск chi-сервера
	r := chi.NewRouter()
	r.Use(corsMiddleware)
	r.Use(middleware.Logger)

	r.Route("/v1", func(r chi.Router) {
		r.Get("/users/admins", getAdminsHandler)
		r.Get("/users/customers", getCustomersHandler)
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	// Логгирование запуска сервера
	fmt.Println("Server starting on :3000")
	http.ListenAndServe(":3000", r)

	log.Printf("server listening at %v", lis.Addr())

	// Обработка ошибки запуска grpc-сервера
	// if err = s.Serve(lis); err != nil {
	// 	log.Fatalf("failed to serve: %v", err)
	// }
}

//
//
//
//
//
//

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

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Пример создания новой записи в БД
// func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
// 	// Делаем запрос на вставку записи в таблицу note
// 	builderInsert := sq.Insert("note").
// 		PlaceholderFormat(sq.Dollar).
// 		Columns("title", "body").
// 		Values(gofakeit.City(), gofakeit.Address().Street).
// 		Suffix("RETURNING id")

// 	query, args, err := builderInsert.ToSql()
// 	if err != nil {
// 		log.Fatalf("failed to build query: %v", err)
// 	}

// 	var noteID int64
// 	err = s.pool.QueryRow(ctx, query, args...).Scan(&noteID)
// 	if err != nil {
// 		log.Fatalf("failed to insert note: %v", err)
// 	}

// 	log.Printf("inserted note with id: %d", noteID)

// 	return &desc.CreateResponse{
// 		Id: noteID,
// 	}, nil
// }
