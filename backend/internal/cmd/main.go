package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
	"wildberies/L0/backend/Kafka/consumer"
	"wildberies/L0/backend/cache"
	"wildberies/L0/backend/internal/app"
	"wildberies/L0/backend/internal/config"
	"wildberies/L0/backend/internal/web/handlers"
	"wildberies/L0/backend/pkg/logger"
	"wildberies/L0/backend/pkg/postgres"

	"context"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	//init environment variables: gobotenv

	//Подгружаем переменные окружения для local из файла .env
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Ошибка загрузки .env файла:", err)
	}

	// init config: cleanenv

	//Инициализация конфига с библиотекой cleanenv (internal/config)
	cfg, err := config.MustLoad()

	if err != nil {
		log.Fatalln(err)
	}

	// init logger: slog

	// Использование логера для вывода log с библиотекой slog
	log := logger.SetupLogger(cfg.Env)

	// Вывод информации в консоль (переменная окружения, версия)
	log.Info(
		"starting url-shortener",
		slog.String("env", cfg.Env),
		slog.String("version", "123"),
	)
	log.Debug("debug messages are enabled")

	//создаёт контекст, который можно отменить вручную
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// init storage: postgre
	pool, err := postgres.NewConn(ctx, cfg)
	if err != nil {
		log.Error("failed to init storage", logger.Err(err))
		os.Exit(1)
	}
	defer pool.Close()

	log.Info("successfully connected to database!")

	// Создание репозитория и кэша
	//newOrder := order.NewOrderRepository(pool)
	cache_data := cache.NewCache()

	orderApp := app.NewApp(pool, log, cache_data)

	// Запуск потребителя Kafka в отдельной горутине, давать app
	go consumer.ConsumerKafka(ctx, orderApp.OrderService, cache_data)

	log.Info("Kafka consumer started")

	// Ожидание сигналов завершения для graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	log.Info("Cache contents", slog.Any("cache", cache_data))

	// Даем время на завершение операций
	time.Sleep(2 * time.Second)
	log.Info("Application shutdown complete")

	// mux := chi.NewRouter()

	// // Добавьте эти строки для обслуживания статических файлов
	//

	// // Ваш существующий API маршрут
	// mux.Get("/order/{id}", handlers.GetOrder(orderApp))

	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	mux.Use(middleware.Logger)

	// API маршрут
	mux.Get("/order/{id}", handlers.GetOrder(orderApp))

	// Обслуживание статических файлов
	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "static"))
	mux.Handle("/*", http.FileServer(filesDir))

	log.Info("Server was started")

	go func() {
		http.ListenAndServe(fmt.Sprintf(":%s", cfg.HTTPServer.Port), mux)
	}()

	// Ожидание сигнала завершения
	sig := <-sigChan
	// реалиация через defer

	fmt.Println()
	log.Info("Received signal, shutting down", slog.String("signal", sig.String()))

	//как работает curl или postmen
	//localhost:8080/order?id=

}
