package config

import (
	"flag"
	"os"
)

type GopherMartConfig struct {
	ServerAddress        string
	Database             string
	AccrualSystemAddress string
}

func NewConfig() *GopherMartConfig {
	return &GopherMartConfig{
		ServerAddress:        "localhost:8080",
		Database:             "postgres://practicum:yandex@localhost:5432/postgres?sslmode=disable",
		AccrualSystemAddress: "http://localhost:8081",
	}
}

func (c *GopherMartConfig) SetConfig() *GopherMartConfig {
	var (
		flagRunAddr          string
		databaseVar          string
		accrualSystemAddress string
	)

	flag.StringVar(&flagRunAddr, "a", c.ServerAddress, "Server address")
	flag.StringVar(&databaseVar, "d", c.Database, "Database connection string")
	flag.StringVar(&accrualSystemAddress, "r", c.AccrualSystemAddress, "Accrual system address")
	flag.Parse()

	if envAddress := os.Getenv("RUN_ADDRESS"); envAddress != "" {
		flagRunAddr = envAddress
	}
	if envDatabase := os.Getenv("DATABASE_URI"); envDatabase != "" {
		databaseVar = envDatabase
	}
	if envAccrualSystemAddress := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envAccrualSystemAddress != "" {
		accrualSystemAddress = envAccrualSystemAddress
	}
	c.ServerAddress = flagRunAddr
	c.Database = databaseVar
	c.AccrualSystemAddress = accrualSystemAddress

	return c
}
