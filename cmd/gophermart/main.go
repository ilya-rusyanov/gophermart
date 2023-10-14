package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/ilya-rusyanov/gophermart/internal/config"
	"github.com/ilya-rusyanov/gophermart/internal/handlers"
	"github.com/ilya-rusyanov/gophermart/internal/httpserver"
	"github.com/ilya-rusyanov/gophermart/internal/logger"
	"github.com/ilya-rusyanov/gophermart/internal/usecases"

	//"github.com/ilya-rusyanov/gophermart/internal/adapters/db"

	"github.com/go-chi/chi"
)

func main() {
	config := config.New()
	config.Parse()

	logger, err := logger.New(config.LogLevel)
	if err != nil {
		fmt.Printf("failed to initialize logger: %q\n", err)
		os.Exit(1)
	}

	context := context.Background()

	//db := db.New(logger, config.DSN)

	registerUsecase := usecases.NewRegister()

	errorHandler := handlers.NewDefaultErrorHandler(logger)

	r := chi.NewRouter()

	r.Route("/api/user", func(r chi.Router) {
		r.Post("/register",
			handlers.NewAuth(
				registerUsecase,
				errorHandler,
			).ServeHTTP)
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
