package postgresql

import "time"

type Config struct {
	Host                  string
	Port                  string
	DbName                string
	UserName              string
	Password              string
	MaxConnections        int32
	MaxConnectionIdleTime time.Duration
}
