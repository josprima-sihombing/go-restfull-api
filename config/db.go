package config

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func InitDB() {
	pool, err := pgxpool.New(context.Background(), Env.DatabaseURL)

	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	if err = pool.Ping(context.Background()); err != nil {
		log.Fatalf("Unable to ping database: %v\n", err)
	}

	DB = pool
	log.Println("✅ Connected to PostgreSQL")
}
