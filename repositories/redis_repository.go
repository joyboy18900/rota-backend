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

// RedisRepository defines the interface for Redis operations
type RedisRepository interface {
	// Token Blacklist methods
	AddToBlacklist(ctx context.Context, token string, ttl time.Duration) error
	IsBlacklisted(ctx context.Context, token string) (bool, error)
}

// redisRepositoryImpl implements RedisRepository
type redisRepositoryImpl struct {
	client *redis.Client
}

// NewRedisRepository creates a new Redis repository
func NewRedisRepository(client *redis.Client) RedisRepository {
	return &redisRepositoryImpl{
		client: client,
	}
}

// Token Blacklist methods
func (r *redisRepositoryImpl) AddToBlacklist(ctx context.Context, token string, ttl time.Duration) error {
	return r.client.Set(ctx, "blacklist:"+token, "1", ttl).Err()
}

func (r *redisRepositoryImpl) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	result, err := r.client.Exists(ctx, "blacklist:"+token).Result()
	return result > 0, err
}

// Future methods for ranking and caching can be added here
// Example:
// func (r *RedisRepository) UpdateRank(ctx context.Context, key string, score float64) error
// func (r *RedisRepository) GetCache(ctx context.Context, key string) (string, error)
