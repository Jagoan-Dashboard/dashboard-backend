
package cache

import (
    "context"
    "fmt"
    "log"
    "building-report-backend/pkg/config"
    
    "github.com/redis/go-redis/v9"
)

func NewRedisClient(cfg config.RedisConfig) *redis.Client {
    client := redis.NewClient(&redis.Options{
        Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
        Password: cfg.Password,
        DB:       cfg.DB,
    })

    ctx := context.Background()
    if err := client.Ping(ctx).Err(); err != nil {
        log.Printf("Warning: Failed to connect to Redis: %v. Some features may be limited.", err)
        // Don't panic, return the client anyway - operations will fail gracefully when used
    }

    return client
}