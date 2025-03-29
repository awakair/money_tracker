package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jonboulle/clockwork"
	"github.com/spf13/cobra"

	"money_tracker/internal/application/services"
	"money_tracker/internal/application/usecases"
	"money_tracker/internal/config"
	httpDelivery "money_tracker/internal/delivery/http"
	"money_tracker/internal/delivery/http/handler"
	"money_tracker/internal/delivery/transaction_rule_worker"
	"money_tracker/internal/infrastructure/cron"
	"money_tracker/internal/infrastructure/id_generator"
	"money_tracker/internal/infrastructure/jwt"
	"money_tracker/internal/infrastructure/password"
	"money_tracker/internal/infrastructure/repository/postgres"
)

// @title Money Tracker API
// @version 1.0
// @description API for managing personal finances, including transactions, categories, and recurring transaction rules.

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	log.SetOutput(os.Stderr)

	cmd := initCmd()
	err := cmd.Execute()
	if err != nil {
		log.Fatalf("App terminated with error: %v", err)
	}
}

func initCmd() *cobra.Command {
	var configPath string

	cmd := &cobra.Command{
		Use:   os.Args[0],
		Short: "Run money tracker service",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.LoadFromFile(configPath)
			if err != nil {
				return fmt.Errorf("can't load config: %w", err)
			}

			return run(cfg)
		},
	}

	cmd.PersistentFlags().StringVarP(&configPath, "config", "c", "config.yaml", "Path to the config yaml file")
	return cmd
}

func run(cfg *config.Config) error {
	logger := initLogger(cfg)
	logger.Info("Starting money tracker service")

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.Database.URL)
	if err != nil {
		return fmt.Errorf("failed to create database pool: %w", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}
	logger.Info("Database connection established")

	db := postgres.NewDatabase(pool)
	clock := clockwork.NewRealClock()
	idGen := id_generator.NewUUIDGenerator()
	passwordHasher := password.NewBcryptHasher()
	jwtService := jwt.NewJWTService(cfg.JWT.Secret, cfg.JWT.Expiration, clock)
	cronService := cron.NewCronService()

	transactionService := services.NewTransactionService(db, clock, cronService, idGen)

	userUseCase := usecases.NewUserUseCase(db, passwordHasher, idGen, clock)
	authUseCase := usecases.NewAuthUseCase(db, passwordHasher, jwtService)
	categoryUseCase := usecases.NewCategoryUseCase(db, idGen, clock)
	transactionUseCase := usecases.NewTransactionUseCase(db, idGen, clock)

	userHandler := handler.NewUserHandler(userUseCase)
	authHandler := handler.NewAuthHandler(authUseCase)
	categoryHandler := handler.NewCategoryHandler(categoryUseCase)
	transactionHandler := handler.NewTransactionHandler(transactionUseCase)
	transactionRuleHandler := handler.NewTransactionRuleHandler(transactionUseCase)

	router := httpDelivery.NewRouter(
		userHandler,
		authHandler,
		transactionHandler,
		categoryHandler,
		transactionRuleHandler,
		jwtService,
		logger,
	)

	server := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router.Setup(),
	}

	worker := transaction_rule_worker.NewWorker(
		transactionService,
		logger,
		clock,
		cfg.Worker.PollInterval,
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	workerErrChan := make(chan error, 1)
	go func() {
		logger.Info("Starting transaction rule worker")
		if err := worker.Run(ctx); err != nil && err != context.Canceled {
			workerErrChan <- fmt.Errorf("worker error: %w", err)
		}
	}()

	serverErrChan := make(chan error, 1)
	go func() {
		logger.Info("Starting HTTP server")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErrChan <- fmt.Errorf("server error: %w", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-sigChan:
		logger.Info("Received shutdown signal", slog.String("signal", sig.String()))
	case err := <-serverErrChan:
		logger.Error("Server error", slog.String("error", err.Error()))
		return err
	case err := <-workerErrChan:
		logger.Error("Worker error", slog.String("error", err.Error()))
		return err
	}

	logger.Info("Shutting down gracefully...")

	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("Server shutdown error", slog.String("error", err.Error()))
		return err
	}

	logger.Info("Service stopped successfully")
	return nil
}

func initLogger(cfg *config.Config) *slog.Logger {
	var logLevel slog.Level
	switch cfg.Log.Level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	})

	return slog.New(handler)
}
