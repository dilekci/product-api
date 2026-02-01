package app

import (
	"product-app/common/postgresql"
	"time"
)

// ConfigurationManager holds application-level configurations.
type ConfigurationManager struct {
	// PostgreSqlConfig contains PostgreSQL related configuration values
	PostgreSqlConfig postgresql.Config
}

// NewConfigurationManager creates and returns a new ConfigurationManager
// with all required configurations initialized.
func NewConfigurationManager() *ConfigurationManager {
	postgreSqlConfig := getPostgreSqlConfig()

	return &ConfigurationManager{
		PostgreSqlConfig: postgreSqlConfig,
	}
}

// getPostgreSqlConfig returns PostgreSQL configuration values.
// In a real-world scenario, these values are typically loaded
// from environment variables or a configuration file.
func getPostgreSqlConfig() postgresql.Config {
	return postgresql.Config{
		Host:                  "localhost",
		Port:                  "6432",
		UserName:              "postgres",
		Password:              "postgres",
		DbName:                "productapp",
		MaxConnections:        10,
		MaxConnectionIdleTime: 30 * time.Second,
	}
}
