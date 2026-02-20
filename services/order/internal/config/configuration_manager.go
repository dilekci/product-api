package config

import (
	"os"
	postgresql "product-app/services/order/internal/adapters/postgresql/common"
	"strconv"
	"time"
)

type ConfigurationManager struct {
	PostgreSqlConfig postgresql.Config
}

func NewConfigurationManager() *ConfigurationManager {
	postgreSqlConfig := getPostgreSqlConfig()

	return &ConfigurationManager{
		PostgreSqlConfig: postgreSqlConfig,
	}
}
func getPostgreSqlConfig() postgresql.Config {
	return postgresql.Config{
		Host:                  getEnvString("DB_HOST", "localhost"),
		Port:                  getEnvString("DB_PORT", "6432"),
		UserName:              getEnvString("DB_USER", "postgres"),
		Password:              getEnvString("DB_PASSWORD", "postgres"),
		DbName:                getEnvString("DB_NAME", "productapp"),
		MaxConnections:        int32(getEnvInt("DB_MAX_CONNECTIONS", 10)),
		MaxConnectionIdleTime: time.Duration(getEnvInt("DB_MAX_IDLE_SECONDS", 30)) * time.Second,
	}
}

func getEnvString(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getEnvInt(key string, fallback int) int {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return parsed
}
