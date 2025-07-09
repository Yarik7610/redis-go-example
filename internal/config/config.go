// internal/config/config.go
package config

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	DB          *sql.DB
	RedisClient *redis.Client
}

func NewConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	db, err := initDB()
	if err != nil {
		log.Fatalf("Init DB error: %v", err)
	}

	redisClient, err := initRedis()
	if err != nil {
		log.Fatalf("Init redis error: %v", err)
	}

	log.Println("CONNECTED TO REDIS SUCCESFULLY")

	return &Config{
		DB:          db,
		RedisClient: redisClient,
	}
}

func initDB() (*sql.DB, error) {
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")
	if dbUser == "" || dbPassword == "" || dbName == "" || dbHost == "" || dbPort == "" {
		return nil, errors.New("Missing required PostgreSQL environment variables")
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func initRedis() (*redis.Client, error) {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	if redisHost == "" || redisPort == "" {
		return nil, errors.New("Missing required Redis environment variables")
	}

	ctx := context.Background()
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: redisPassword,
		DB:       0,
	})
	if err := redisClient.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return redisClient, nil
}

func (c *Config) Close() {
	if c.DB != nil {
		c.DB.Close()
	}
	if c.RedisClient != nil {
		c.RedisClient.Close()
	}
}
