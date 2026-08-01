package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"screw-it-desk-backend/services/users/internal/closer"
	"screw-it-desk-backend/services/users/internal/config"
	handlers "screw-it-desk-backend/services/users/internal/handlers/users"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

type App struct {
	serviceProvider *serviceProvider
	httpServer      *chi.Mux
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	return a.runHTTPServer()
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initHTTPServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	err := config.Load(".env")
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	// Инициализация HTTP Сервера и роутов
	a.httpServer = chi.NewRouter()

	// TODO найти как подключать cors-middleware
	// Add middleware

	a.httpServer.Use(corsMiddleware)
	a.httpServer.Use(middleware.Logger)
	a.httpServer.Use(middleware.Recoverer)
	a.httpServer.Use(middleware.RequestID)

	a.httpServer.Route("/v1", func(r chi.Router) {
		r.Get("/users/admins", handlers.GetAdminsHandler)
		r.Get("/users/customers", handlers.GetCustomersHandler)
	})

	a.httpServer.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	// Health check endpoint
	// @Summary Health Check
	// @Description Returns the health status of the payment service
	// @Tags Health
	// @Produce json
	// @Success 200 {object} map[string]string "Service is healthy"
	// @Router /health [get]
	a.httpServer.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"status":"ok","service":"user-service"}`)
	})

	// TODO проверить путь для swagger-документации и исправить на относительный
	a.httpServer.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	return nil
}

func (a *App) runHTTPServer() error {
	log.Print("a.serviceProvider.HTTPConfig().Address() ===  " + a.serviceProvider.HTTPConfig().Address())

	err := http.ListenAndServe(":3000", a.httpServer)

	if err != nil {
		return err
	}
	// Логгирование запуска сервера
	log.Printf("HTTP server is running on %s", a.serviceProvider.HTTPConfig().Address())

	return nil
}

// TODO вынести инициализацию middleware
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// w.Header().Set("Access-Control-Allow-Origin", "*")
		// w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		// w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		// w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
