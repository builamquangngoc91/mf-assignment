package configs

import "os"

type Database struct {
	Host     string
	Username string
	Password string
	Name     string
	Port     string
}

type BankingService struct {
	Port string
}

type Config struct {
	Database       Database
	BankingService BankingService
}

var Cfg Config

func LoadConfig() {
	Cfg = Config{
		Database: Database{
			Host:     os.Getenv("BANKING_DB_HOST"),
			Username: os.Getenv("BANKING_DB_USERNAME"),
			Password: os.Getenv("BANKING_DB_PASSWORD"),
			Name:     os.Getenv("BANKING_DB_NAME"),
			Port:     os.Getenv("BANKING_DB_PORT"),
		},
		BankingService: BankingService{
			Port: os.Getenv("BANKING_SERVICE_PORT"),
		},
	}
}
