package config

import (
	"flag"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config is a struct for configuration
type GopherMartConfig struct {
	ServerAddress        string `yaml:"host"`
	Database             string `yaml:"dblink"`
	AccrualSystemAddress string `yaml:"accrualsystemaddr"`
	LogLevel             string `yaml:"loglevel"`
	Key                  string `yaml:"key"`
}

// NewConfig is a constructor for Config
func NewConfig() *GopherMartConfig {
	return &GopherMartConfig{
		ServerAddress:        "localhost:8080",
		Database:             "postgres://practicum:yandex@localhost:5432/postgres?sslmode=disable",
		AccrualSystemAddress: "http://localhost:8081",
		LogLevel:             "info",
	}
}

// SetConfig is a function for setting configuration
//func (c *GopherMartConfig) SetConfig() *GopherMartConfig {
//	var (
//		flagRunAddr          string
//		databaseVar          string
//		accrualSystemAddress string
//	)
//	//handle flags
//	flag.StringVar(&flagRunAddr, "a", c.ServerAddress, "Server address")
//	flag.StringVar(&databaseVar, "d", c.Database, "Database connection string")
//	flag.StringVar(&accrualSystemAddress, "r", c.AccrualSystemAddress, "Accrual system address")
//	flag.Parse()
//	//handle env variables
//	if envAddress := os.Getenv("RUN_ADDRESS"); envAddress != "" {
//		flagRunAddr = envAddress
//	}
//	if envDatabase := os.Getenv("DATABASE_URI"); envDatabase != "" {
//		databaseVar = envDatabase
//	}
//	if envAccrualSystemAddress := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envAccrualSystemAddress != "" {
//		accrualSystemAddress = envAccrualSystemAddress
//	}
//	//set config
//	c.ServerAddress = flagRunAddr
//	c.Database = databaseVar
//	c.AccrualSystemAddress = accrualSystemAddress
//
//	return c
//}

func GetConfig() *GopherMartConfig {
	var config *GopherMartConfig
	file, err := os.Open("./config/config.yaml")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println(err)
	}
	var (
		flagRunAddr          string
		databaseVar          string
		accrualSystemAddress string
	)
	//handle flags
	flag.StringVar(&flagRunAddr, "a", config.ServerAddress, "Server address")
	flag.StringVar(&databaseVar, "d", config.Database, "Database connection string")
	flag.StringVar(&accrualSystemAddress, "r", config.AccrualSystemAddress, "Accrual system address")
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
	config.ServerAddress = flagRunAddr
	config.Database = databaseVar
	config.AccrualSystemAddress = accrualSystemAddress

	return config
}
