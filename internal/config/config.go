package config

import (
	"flag"
	"os"
)

type Config struct {
	ListenAddr  string
	DSN         string
	AccrualAddr string
}

func New() *Config {
	var res Config
	flag.StringVar(&res.ListenAddr, "a", ":8080", "address and port to listen on")
	flag.StringVar(&res.DSN, "d", "host=localhost port=5433 user=gophermart dbname=gophermart sslmode=disable", "database URI")
	flag.StringVar(&res.AccrualAddr, "r", "localhost:8081", "accrual system address")
	return &res
}

func (c *Config) Parse() {
	flag.Parse()

	if val := os.Getenv("RUN_ADDRESS"); len(val) > 0 {
		c.ListenAddr = val
	}

	if val := os.Getenv("DATABASE_URI"); len(val) > 0 {
		c.DSN = val
	}

	if val := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); len(val) > 0 {
		c.AccrualAddr = val
	}
}
