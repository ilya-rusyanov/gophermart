package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/ilya-rusyanov/gophermart/internal/adapters/accrual"
	"github.com/ilya-rusyanov/gophermart/internal/adapters/httpserver"
	"github.com/ilya-rusyanov/gophermart/internal/adapters/httpserver/middleware"
	"github.com/ilya-rusyanov/gophermart/internal/adapters/postgres"
	"github.com/ilya-rusyanov/gophermart/internal/config"
	"github.com/ilya-rusyanov/gophermart/internal/handlers"
	"github.com/ilya-rusyanov/gophermart/internal/logger"
	"github.com/ilya-rusyanov/gophermart/internal/storage"
	"github.com/ilya-rusyanov/gophermart/internal/usecases"

	"github.com/go-chi/chi"
)

const tokenExpiration = time.Hour * 24 * 7
const signingKey = "TODO"

func main() {
	config := config.New()
	err := config.Parse()
	if err != nil {
		fmt.Printf("failed to parse configuration: %q\n", err)
		os.Exit(1)
	}

	logger, err := logger.New(config.LogLevel)
	if err != nil {
		fmt.Printf("failed to initialize logger: %q\n", err)
		os.Exit(1)
	}

	context, cancel := context.WithCancel(context.Background())
	defer cancel()

	db := postgres.MustInit(context, logger, config.DSN, config.MaxUserNameLen)
	defer db.Close()

	accrualAdapter := accrual.New(logger, config.AccrualAddr)

	userStorage := storage.NewUser(db)
	orderStorage := storage.NewOrder(db, logger)
	accrualStorage := storage.NewAccrual(db)
	balanceStorage := storage.NewBalance(db)
	withdrawalStorage := storage.NewWithdrawal(db)

	registerUsecase := usecases.NewRegister(
		logger, tokenExpiration, signingKey, userStorage,
	)
	loginUsecase := usecases.NewLogin(
		logger, tokenExpiration, signingKey, userStorage,
	)
	createOrderUsecase := usecases.NewCreateOrder(
		logger, orderStorage,
	)

	feedAccrual := usecases.NewFeedAccrual(
		logger, accrualStorage, accrualAdapter)
	go feedAccrual.Run(context, 1*time.Second)

	withdrawUsecase := usecases.NewWithdraw(withdrawalStorage)

	errorHandler := handlers.NewDefaultErrorHandler(logger).Handle

	r := chi.NewRouter()

	r.Use(middleware.NewLogger(logger).Middleware)
	r.Use(middleware.NewAuth(
		signingKey,
		errorHandler,
		"/api/user/register",
		"/api/user/login",
	).Middleware)

	r.Route("/api/user", func(r chi.Router) {
		r.Post("/register", handlers.NewAuth(registerUsecase, errorHandler).ServeHTTP)
		r.Post("/login", handlers.NewAuth(loginUsecase, errorHandler).ServeHTTP)
		r.Post("/orders", handlers.NewOrderCreation(logger, createOrderUsecase, errorHandler).ServeHTTP)
		r.Get("/orders", handlers.NewListOrders(logger, orderStorage, errorHandler).ServeHTTP)
		r.Route("/balance", func(r chi.Router) {
			r.Get("/", handlers.NewShowBalance(logger, balanceStorage, errorHandler).ServeHTTP)
			r.Post("/withdraw", handlers.NewWithdraw(logger, withdrawUsecase, errorHandler).ServeHTTP)
		})
		r.Get("/withdrawals", handlers.NewListWithdrawals(logger, withdrawalStorage, errorHandler).ServeHTTP)
	})

	httpServer := httpserver.New(config.ListenAddr, r)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	select {
	case <-interrupt:
		logger.Info("interrupt")
	case err = <-httpServer.Error():
		logger.Errorf("http server error: %q", err)
	}

	err = httpServer.Shutdown(context)
	if err != nil {
		logger.Error(err)
	}
}
