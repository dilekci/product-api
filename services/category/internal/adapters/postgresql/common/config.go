package postgresql

import "time"

// Config holds PostgreSQL database connection and pool configuration.
// It is intended to be used for initializing and managing a PostgreSQL client.
type Config struct {
	Host                  string
	Port                  string
	DbName                string
	UserName              string
	Password              string
	MaxConnections        int32
	MaxConnectionIdleTime time.Duration
}
