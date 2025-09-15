
package redis

import (
    "context"
    "encoding/json"
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
    return r.client.Set(ctx, key, data, expiration).Err()
}

func (r *cacheRepositoryImpl) Get(ctx context.Context, key string, dest interface{}) error {
    data, err := r.client.Get(ctx, key).Result()
    if err != nil {
        return err
    }
    return json.Unmarshal([]byte(data), dest)
}

func (r *cacheRepositoryImpl) Delete(ctx context.Context, keys ...string) error {
    return r.client.Del(ctx, keys...).Err()
}

func (r *cacheRepositoryImpl) Exists(ctx context.Context, key string) (bool, error) {
    result, err := r.client.Exists(ctx, key).Result()
    if err != nil {
        return false, err
    }
    return result > 0, nil
}

func (r *cacheRepositoryImpl) Flush(ctx context.Context) error {
    return r.client.FlushAll(ctx).Err()
}