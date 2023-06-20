package config

import (
	"flag"
	"os"
)

// Config is a struct for configuration
type GopherMartConfig struct {
	ServerAddress        string
	Database             string
	AccrualSystemAddress string
}

// NewConfig is a constructor for Config
func NewConfig() *GopherMartConfig {
	return &GopherMartConfig{
		ServerAddress:        "localhost:8080",
		Database:             "postgres://practicum:yandex@localhost:5432/postgres?sslmode=disable",
		AccrualSystemAddress: "http://localhost:8081",
	}
}

// SetConfig is a function for setting configuration
func (c *GopherMartConfig) SetConfig() *GopherMartConfig {
	var (
		flagRunAddr          string
		databaseVar          string
		accrualSystemAddress string
	)
	//handle flags
	flag.StringVar(&flagRunAddr, "a", c.ServerAddress, "Server address")
	flag.StringVar(&databaseVar, "d", c.Database, "Database connection string")
	flag.StringVar(&accrualSystemAddress, "r", c.AccrualSystemAddress, "Accrual system address")
	flag.Parse()
	//handle env variables
	if envAddress := os.Getenv("RUN_ADDRESS"); envAddress != "" {
		flagRunAddr = envAddress
	}
	if envDatabase := os.Getenv("DATABASE_URI"); envDatabase != "" {
		databaseVar = envDatabase
	}
	if envAccrualSystemAddress := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envAccrualSystemAddress != "" {
		accrualSystemAddress = envAccrualSystemAddress
	}
	//set config
	c.ServerAddress = flagRunAddr
	c.Database = databaseVar
	c.AccrualSystemAddress = accrualSystemAddress

	return c
}
