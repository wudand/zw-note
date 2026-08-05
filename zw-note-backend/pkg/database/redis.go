package database

// Redis support is reserved for future use (caching, token blacklist, rate limiting, etc.).
// Uncomment and implement when needed using github.com/redis/go-redis/v9.
//
// import (
// 	"context"
// 	"fmt"
// 	"github.com/redis/go-redis/v9"
// )
//
// type RedisConfig struct {
// 	Addr     string
// 	Password string
// 	DB       int
// }
//
// func NewRedis(cfg RedisConfig) (*redis.Client, error) {
// 	rdb := redis.NewClient(&redis.Options{
// 		Addr:     cfg.Addr,
// 		Password: cfg.Password,
// 		DB:       cfg.DB,
// 	})
// 	if err := rdb.Ping(context.Background()).Err(); err != nil {
// 		return nil, fmt.Errorf("ping redis: %w", err)
// 	}
// 	return rdb, nil
// }
