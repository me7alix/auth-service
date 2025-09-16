package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

type DBClient struct {
	Client *sql.DB
}

func NewDBClient(conninfo string) *DBClient {
	db, err := sql.Open("postgres", conninfo)
	if err != nil {
		fmt.Println("DB ERROR")
		log.Fatal(err)
	}

	return &DBClient{Client: db}
}

func NewRedisClient(conninfo string) *redis.Client {
	opts, err := redis.ParseURL(conninfo)
	if err != nil {
		fmt.Println("REDIS ERROR")
		log.Fatal(err)
	}

	return redis.NewClient(opts)
}
