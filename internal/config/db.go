package config

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/chandraRGB/MiniProject-GustiChandraMiftahulMunir/internal/domain"
)

// DBConfig holds database configuration values.
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

// getDefaultDBConfig returns default DB config and allows override by environment variables.
func getDefaultDBConfig() DBConfig {
	cfg := DBConfig{
		Host:     "127.0.0.1",
		Port:     "3306",
		User:     "root",
		Password: "",
		Name:     "evermos",
	}

	if v := os.Getenv("DB_HOST"); v != "" {
		cfg.Host = v
	}
	if v := os.Getenv("DB_PORT"); v != "" {
		cfg.Port = v
	}
	if v := os.Getenv("DB_USER"); v != "" {
		cfg.User = v
	}
	if v := os.Getenv("DB_PASSWORD"); v != "" {
		cfg.Password = v
	}
	if v := os.Getenv("DB_NAME"); v != "" {
		cfg.Name = v
	}

	return cfg
}

// NewDB initializes a new GORM MySQL connection and runs migrations.
func NewDB() (*gorm.DB, error) {
	cfg := getDefaultDBConfig()

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto-migrate all domain models
	if err := db.AutoMigrate(
		&domain.User{},
		&domain.Toko{},
		&domain.Alamat{},
		&domain.Category{},
		&domain.Produk{},
		&domain.FotoProduk{},
		&domain.Trx{},
		&domain.LogProduk{},
		&domain.DetailTrx{},
	); err != nil {
		return nil, err
	}

	return db, nil
}
