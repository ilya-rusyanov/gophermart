package main

import (
	"github.com/ilya-rusyanov/gophermart/internal/config"
)

func main() {
	config := config.New()
	config.Parse()
}
