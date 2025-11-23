package utils

import (
	"time"

	"github.com/ADMex1/GoProject/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Generate Token Generate Refresh Token
func GenerateToken(userID int64, role string, email string, publicID uuid.UUID) (string, error) {
	secret := config.AppConfig.JWTSecret
	duration, _ := time.ParseDuration(config.AppConfig.JWTExpire)

	claims := jwt.MapClaims{
		"user_id":  userID,
		"role":     role,
		"pub_id":   publicID,
		"email":    email,
		"exp_time": time.Now().Add(duration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}
func RefreshToken(userID int64) (string, error) {
	secret := config.AppConfig.JWTSecret
	duration, _ := time.ParseDuration(config.AppConfig.JWTRefreshToken)
	claims := jwt.MapClaims{
		"user_id":  userID,
		"exp_time": time.Now().Add(duration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
