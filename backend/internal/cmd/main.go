package main

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"wildberies/L0/backend/Kafka/consumer"
	"wildberies/L0/backend/cache"
	"wildberies/L0/backend/internal/app"
	"wildberies/L0/backend/internal/config"
	router "wildberies/L0/backend/internal/web"
	"wildberies/L0/backend/pkg/logger"
	"wildberies/L0/backend/pkg/postgres"

	"context"

	"github.com/joho/godotenv"
)

func main() {

	//1. Инициализация переменных окружения

	if err := godotenv.Load(); err != nil {
		if !os.IsNotExist(err) {
			log.Fatal("Ошибка загрузки .env файла:", err)
		}
	}

	//2. Инициализация config
	cfg, err := config.Load()

	if err != nil {
		log.Fatalln(err)
	}

	//3. Создание logger
	log := logger.NewLogger(cfg.Env)

	//4. Вывод информации в консоль (переменная окружения, версия)
	log.Info(
		"starting project",
		slog.String("env", cfg.Env),
		slog.String("version", "123"),
	)
	log.Debug("debug messages are enabled")

	if err := run(*cfg, log); err != nil {
		log.Error("application stopped with error", logger.Err(err))
		os.Exit(1)
	}

}

func run(cfg config.Config, log *slog.Logger) error {
	//Общий context приложения
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	//Создания подключения к PostgreSQL
	pool, err := postgres.NewConn(ctx, &cfg.PostgresConfig)

	if err != nil {
		log.Error("failed to init storage", logger.Err(err))
		return err
	}
	defer pool.Close()

	log.Info("successfully connected to database!")

	//Создание Cache
	cacheData := cache.NewCache(cfg.MaxItemsCache)

	//Создание App
	orderApp := app.NewApp(pool, cacheData, log)
	err = orderApp.OrderService.WarmCache(ctx, cfg.LimitCache)

	if err != nil {
		log.Error("failed to warm cache", logger.Err(err))
		return err
	}

	//Создание router
	rout, err := router.BuildRouter(orderApp.OrderService)

	if err != nil {
		log.Error("get workdir", logger.Err(err))
		return err
	}

	var errCons error
	consumerErr := make(chan error, 1)
	//Создание Kafka consumer
	go func() {
		consumerErr <- consumer.ConsumerKafka(ctx, orderApp.OrderService, cfg, log)
	}()

	//Создание HTTP server
	srv := router.BuildServer(cfg, rout)

	srvErr := make(chan error, 1)

	go func() {
		log.Info("http server started", slog.String("addr", srv.Addr))

		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("failed to init http-server", logger.Err(err))
			srvErr <- err
		}
	}()

	select {
	case <-ctx.Done():
		log.Info("shutdown signal received", slog.Any("reason", ctx.Err()))
	case err := <-consumerErr:
		if err != nil {
			log.Error("failed consumer kafka:", logger.Err(err))
			errCons = err
			cancel()
		}
	case err := <-srvErr:
		log.Error("failed server connection:", logger.Err(err))
		errCons = err
		cancel()
	}

	//Создание shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("shutdown server")
		return fmt.Errorf("shutdown server: %w", err)
	}

	log.Info("application stopped")
	if errCons != nil {
		return errCons
	}
	return nil
}
