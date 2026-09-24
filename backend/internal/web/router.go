package router

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"wildberies/L0/backend/internal/config"
	"wildberies/L0/backend/internal/web/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func BuildRouter(orderApp handlers.OrderGetter) (http.Handler, error) {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	mux.Use(middleware.Logger)

	// API маршрут
	mux.Get("/order/{id}", handlers.GetOrder(orderApp))

	workDir, err := os.Getwd()

	if err != nil {
		return nil, err
	}

	filesDir := http.Dir(filepath.Join(workDir, "static"))
	mux.Handle("/*", http.FileServer(filesDir))
	return mux, nil
}

func BuildServer(cfg config.Config, router http.Handler) *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.HTTPServer.Port),
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}
}
