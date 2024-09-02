package services

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type ICacheService interface {
	Delete(key string) error
	BlacklistTkn(token string, expiration time.Duration) error
	IsTknBlacklisted(token string) (bool, error)
}

type cacheService struct {
	client *redis.Client
}

func NewCacheService(redisAddr, redisPassword string, redisDB int) ICacheService {
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})

	return &cacheService{
		client: rdb,
	}
}

func (cs *cacheService) BlacklistTkn(token string, expiration time.Duration) error {
    return cs.client.Set(context.Background(), token, "blacklisted", expiration).Err()
}

func (cs *cacheService) IsTknBlacklisted(token string) (bool, error) {
    result, err := cs.client.Get(context.Background(), token).Result()
    if err == redis.Nil {
        return false, nil 
    } else if err != nil {
        return false, err
    }
    return result == "blacklisted", nil
}

func (cs *cacheService) Delete(key string) error {
	err := cs.client.Del(context.Background(), key).Err()
	if err != nil {
		return err
	}
	return nil
}
