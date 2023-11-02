package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	ListenAddr     string
	DSN            string
	AccrualAddr    string
	LogLevel       string
	MaxUserNameLen int
}

func New() *Config {
	var res Config
	flag.StringVar(&res.ListenAddr, "a", ":8080", "address and port to listen on")
	flag.StringVar(&res.DSN, "d", "host=localhost port=5433 user=gophermart dbname=gophermart sslmode=disable", "database URI")
	flag.StringVar(&res.AccrualAddr, "r", "http://localhost:8082", "accrual system address")
	flag.StringVar(&res.LogLevel, "l", "info", "log level")
	flag.IntVar(&res.MaxUserNameLen, "m", 40, "maximum username length")
	return &res
}

func (c *Config) Parse() error {
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

	if val := os.Getenv("LOG_LEVEL"); len(val) > 0 {
		c.LogLevel = val
	}

	if val := os.Getenv("MAX_USERNAME_LEN"); len(val) > 0 {
		length, err := strconv.Atoi(val)
		if err != nil {
			return fmt.Errorf("cannot parse username length: %w", err)
		}

		c.MaxUserNameLen = length
	}

	return nil
}
