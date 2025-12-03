package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB        *gorm.DB
	AppConfig *Config
)

type Config struct {
	AppPort          string
	DBHost           string
	DBUser           string
	DBPassword       string
	DBName           string
	DBPort           string
	JWTSecret        string
	JWTExpireMinutes string
	JWTRefreshToken  string
	JWTExpire        string
	APPURL           string
}

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Not connected to the .env file")
	}
	AppConfig = &Config{
		AppPort:          getEnv("PORT", "3000"),
		DBHost:           getEnv("DB_HOST", "localhost"),
		DBName:           getEnv("DB_NAME", "ProjectManagement"),
		DBUser:           getEnv("DB_USER", "postgres"),
		DBPassword:       getEnv("DB_PASSWORD", "NeuXana"),
		DBPort:           getEnv("DB_PORT", "5432"),
		JWTSecret:        getEnv("JWT_SECRET", "supersecret"),
		JWTExpire:        getEnv("JWT_EXPIRATION", "24h"),
		JWTExpireMinutes: getEnv("JWT_EXPIRY_MINUTES", "60"),
		JWTRefreshToken:  getEnv("REFRESH_TOKEN_EXPIRED", "24h"),
		APPURL:           getEnv("APP_URL", "http://localhost:3000"),
	}
}

func getEnv(key string, fallback string) string {
	values, exist := os.LookupEnv(key)

	if exist {
		return values
	} else {
		return fallback
	}
}

func DBConnect() {
	cfg := AppConfig

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("DB Connection Failure", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Getting database Instance failed", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db
}
