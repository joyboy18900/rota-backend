package repositories

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisConfig represents the configuration for Redis connection
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// RedisRepository handles Redis operations
type RedisRepository struct {
	client *redis.Client
}

// NewRedisRepository creates a new Redis repository
func NewRedisRepository(client *redis.Client) *RedisRepository {
	return &RedisRepository{
		client: client,
	}
}

// Token Blacklist methods
func (r *RedisRepository) AddToBlacklist(ctx context.Context, token string, ttl time.Duration) error {
	return r.client.Set(ctx, "blacklist:"+token, "1", ttl).Err()
}

func (r *RedisRepository) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	result, err := r.client.Exists(ctx, "blacklist:"+token).Result()
	return result > 0, err
}

// Future methods for ranking and caching can be added here
// Example:
// func (r *RedisRepository) UpdateRank(ctx context.Context, key string, score float64) error
// func (r *RedisRepository) GetCache(ctx context.Context, key string) (string, error)
