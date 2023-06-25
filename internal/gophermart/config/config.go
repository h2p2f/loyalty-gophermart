package config

import (
	"flag"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type GopherMartConf struct {
	Server struct {
		Address string `yaml:"host"`
	} `yaml:"server"`
	Database struct {
		Database string `yaml:"dblink"`
	} `yaml:"database"`
	AccrualSystem struct {
		Address string `yaml:"accrualsystemaddr"`
	} `yaml:"accrualsystem"`
	SystemParams struct {
		LogLevel string `yaml:"loglevel"`
		Key      string `yaml:"key"`
	} `yaml:"systemparams"`
}

// GetConfig get config from file and envs

func GetConfig() *GopherMartConf {
	var config *GopherMartConf
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
	flag.StringVar(&flagRunAddr, "a", config.Server.Address, "Server address")
	flag.StringVar(&databaseVar, "d", config.Database.Database, "Database connection string")
	flag.StringVar(&accrualSystemAddress, "r", config.AccrualSystem.Address, "Accrual system address")
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
	config.Server.Address = flagRunAddr
	config.Database.Database = databaseVar
	config.AccrualSystem.Address = accrualSystemAddress

	return config
}
