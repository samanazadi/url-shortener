package database

import (
	"fmt"
	"github.com/samanazadi/url-shortener/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewGormDB(cfg *config.Config) (*gorm.DB, error) {
	dbuser := cfg.DBUser
	dbpass := cfg.DBPass
	dbhost := cfg.DBHost
	dbport := cfg.DBPort
	dbdb := cfg.DBName
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbhost, dbuser, dbpass, dbdb, dbport)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
