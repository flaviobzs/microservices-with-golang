package db

import (
	"github.com/go-redis/redis/v7"
	"os"
)
// importar pacode do client redis 

func Connect() *redis.Client {
	// conectar com o redis 
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: "",
		DB:       0,
	})
	return client
}
