package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ilya-rusyanov/gophermart/internal/config"
	"github.com/ilya-rusyanov/gophermart/internal/logger"

	//"github.com/ilya-rusyanov/gophermart/internal/adapters/db"
	ht "github.com/ilya-rusyanov/gophermart/internal/adapters/http"
	"github.com/ilya-rusyanov/gophermart/internal/usecases"

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

	//db := db.New(logger, config.DSN)
	authUsecase := usecases.NewAuth( /*&db*/ )
	httpAdapter := ht.New(logger, authUsecase)

	r := chi.NewRouter()

	r.Route("/api/user", func(r chi.Router) {
		r.Post("/register", httpAdapter.Register)
	})

	http.ListenAndServe(config.ListenAddr, r)
}
