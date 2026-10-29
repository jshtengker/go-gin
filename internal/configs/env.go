package configs

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv                         string
	ServerPort                     string
	DBHost                         string
	DBPort                         string
	DBUser                         string
	DBPassword                     string
	DBName                         string
	JWTSecret                      string
	JWTExpiryMinutes               int
	RefreshTokenSecret             string
	RefreshTokenExpiryHours        int
	RefreshTokenGracePeriodSeconds int
}

func LoadEnv() {
	_ = godotenv.Load("internal/configs/.env")
}

func Load() *Config {
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "development"
	}

	jwtExpiry := 15
	if val, err := strconv.Atoi(os.Getenv("JWT_EXPIRY_MINUTES")); err == nil && val > 0 {
		jwtExpiry = val
	}

	refreshExpiry := 168
	if val, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_EXPIRY_HOURS")); err == nil && val > 0 {
		refreshExpiry = val
	}

	gracePeriod := 30
	if val, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_GRACE_PERIOD_SECONDS")); err == nil && val >= 0 {
		gracePeriod = val
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default_jwt_secret_key_change_me"
	}

	refreshTokenSecret := os.Getenv("REFRESH_TOKEN_SECRET")
	if refreshTokenSecret == "" {
		refreshTokenSecret = "default_refresh_token_secret_key_change_me"
	}

	return &Config{
		AppEnv:                         appEnv,
		ServerPort:                     os.Getenv("SERVER_PORT"),
		DBHost:                         os.Getenv("DB_HOST"),
		DBPort:                         os.Getenv("DB_PORT"),
		DBUser:                         os.Getenv("DB_USER"),
		DBPassword:                     os.Getenv("DB_PASSWORD"),
		DBName:                         os.Getenv("DB_NAME"),
		JWTSecret:                      jwtSecret,
		JWTExpiryMinutes:               jwtExpiry,
		RefreshTokenSecret:             refreshTokenSecret,
		RefreshTokenExpiryHours:        refreshExpiry,
		RefreshTokenGracePeriodSeconds: gracePeriod,
	}
}
