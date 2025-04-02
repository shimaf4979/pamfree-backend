// backend/config/config.go
package config

import (
	"os"
	"fmt"

	"github.com/joho/godotenv"
)

// Config はアプリケーション設定を格納する構造体
type Config struct {
	Env             string
	ServerAddress   string
	DBHost          string
	DBPort          string
	DBUser          string
	DBPassword      string
	DBName          string
	JWTSecret       string
	CloudinaryName  string
	CloudinaryKey   string
	CloudinarySecret string
}

// LoadConfig は環境変数から設定を読み込む
func LoadConfig() (*Config, error) {
	// 開発環境では.envファイルを読み込む
	if os.Getenv("ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			return nil, fmt.Errorf(".envファイルの読み込みに失敗しました: %w", err)
		}
	}

	config := &Config{
		Env:             getEnv("ENV", "development"),
		ServerAddress:   getEnv("SERVER_ADDRESS", ":8080"),
		DBHost:          getEnv("DB_HOST", "localhost"),
		DBPort:          getEnv("DB_PORT", "3306"),
		DBUser:          getEnv("DB_USER", "root"),
		DBPassword:      getEnv("DB_PASSWORD", "password"),
		DBName:          getEnv("DB_NAME", "mapapp"),
		JWTSecret:       getEnv("JWT_SECRET", "your-secret-key"),
		CloudinaryName:  getEnv("CLOUDINARY_CLOUD_NAME", ""),
		CloudinaryKey:   getEnv("CLOUDINARY_API_KEY", ""),
		CloudinarySecret: getEnv("CLOUDINARY_API_SECRET", ""),
	}

	return config, nil
}

// getEnv は環境変数を取得し、未設定の場合はデフォルト値を返す
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}