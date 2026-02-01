package postgresql

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

func GetConnectionPool(ctx context.Context, config Config) *pgxpool.Pool {

	connString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host,
		config.Port,
		config.UserName,
		config.Password,
		config.DbName,
	)

	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		panic(err)
	}

	// defaults
	if config.MaxConnections > 0 {
		poolConfig.MaxConns = config.MaxConnections
	} else {
		poolConfig.MaxConns = 5
	}

	if config.MaxConnectionIdleTime > 0 {
		poolConfig.MaxConnIdleTime = config.MaxConnectionIdleTime
	} else {
		poolConfig.MaxConnIdleTime = 30 * time.Second
	}

	pool, err := pgxpool.ConnectConfig(ctx, poolConfig)
	if err != nil {
		panic(err)
	}

	return pool
}
