package datastore

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"

	"worknote-api/config"
)

var (
	// DB is the PostgreSQL database connection
	DB *sqlx.DB
	// Redis is the Redis client
	Redis *redis.Client
)

// Initialize sets up all data store connections (PostgreSQL, Redis)
func Initialize() {
	initPostgres()
	initRedis()
}

// Close closes all data store connections
func Close() {
	if DB != nil {
		DB.Close()
	}
	if Redis != nil {
		Redis.Close()
	}
}

func initPostgres() {
	cfg := config.Get()

	var err error
	DB, err = sqlx.Connect("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	log.Info("PostgreSQL connected")
}

func initRedis() {
	cfg := config.Get()

	Redis = redis.NewClient(&redis.Options{
		Addr:     cfg.RedisURL,
		Password: cfg.RedisPassword,
		DB:       0,
	})

	log.Info("Redis connected")
}
