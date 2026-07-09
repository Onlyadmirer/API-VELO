package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis() *redis.Client {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		log.Println("REDIS_URL belum di-set di .env")
	}

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Println("gagal membaca Redis url: ", err)
	}

	client := redis.NewClient(opt)

	ctx := context.Background()
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		log.Println("gagal konek dengan redis upstash")
	} else {

		fmt.Println("Berhasil konek ke redis upstash!:", pong)
	}

	return client
}
