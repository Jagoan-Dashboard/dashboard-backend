
package redis

import (
    "context"
    "encoding/json"
    "log"
    "time"
    
    "building-report-backend/internal/domain/repository"
    "github.com/redis/go-redis/v9"
)

type cacheRepositoryImpl struct {
    client *redis.Client
}

func NewCacheRepository(client *redis.Client) repository.CacheRepository {
    return &cacheRepositoryImpl{client: client}
}

func (r *cacheRepositoryImpl) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
    data, err := json.Marshal(value)
    if err != nil {
        return err
    }
    
    if err := r.client.Set(ctx, key, data, expiration).Err(); err != nil {
        log.Printf("Warning: Failed to set cache key %s: %v", key, err)
        // Don't return error, just log it - caching should be optional
    }
    return nil
}

func (r *cacheRepositoryImpl) Get(ctx context.Context, key string, dest interface{}) error {
    data, err := r.client.Get(ctx, key).Result()
    if err != nil {
        if err == redis.Nil {
            // Key does not exist, return the error so the caller can handle it appropriately
            return err
        }
        log.Printf("Warning: Failed to get cache key %s: %v", key, err)
        // Return error so that the caller knows to fetch from the database
        return err
    }
    
    return json.Unmarshal([]byte(data), dest)
}

func (r *cacheRepositoryImpl) Delete(ctx context.Context, keys ...string) error {
    if err := r.client.Del(ctx, keys...).Err(); err != nil {
        log.Printf("Warning: Failed to delete cache keys %v: %v", keys, err)
        // Don't return error, just log it - deletion failure shouldn't break functionality
    }
    return nil
}

func (r *cacheRepositoryImpl) Exists(ctx context.Context, key string) (bool, error) {
    result, err := r.client.Exists(ctx, key).Result()
    if err != nil {
        log.Printf("Warning: Failed to check existence of cache key %s: %v", key, err)
        // Return false and nil error - assume key doesn't exist if Redis is unavailable
        return false, nil
    }
    return result > 0, nil
}

func (r *cacheRepositoryImpl) Flush(ctx context.Context) error {
    if err := r.client.FlushAll(ctx).Err(); err != nil {
        log.Printf("Warning: Failed to flush Redis: %v", err)
        // Don't return error, just log it - flushing should be optional
    }
    return nil
}