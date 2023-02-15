package redis

import (
	"context"

	redisClient "example.com/GoService/pkg/redis"

	"github.com/rs/zerolog"
)

type CacheClient interface {
	CachedData(ctx context.Context, data string, uID string) error
	GetDataByID(ctx context.Context, uID string) (string, error)
	DeleteDataByID(ctx context.Context, uID string) error
}

type cacheClient struct {
	RedisConnection redisClient.Client
	logger          zerolog.Logger
}

func NewCacheClient(conn redisClient.Client, log zerolog.Logger) CacheClient {
	return &cacheClient{
		RedisConnection: conn,
		logger:          log,
	}
}

func (sr *cacheClient) CachedData(ctx context.Context, data string, uID string) error {
	err := sr.RedisConnection.Set(uID, data)
	if err != nil {
		return err
	}
	return nil
}

func (sr *cacheClient) GetDataByID(ctx context.Context, uID string) (string, error) {
	data, err := sr.RedisConnection.GetValue(uID)
	if err != nil {
		sr.logger.Error().Err(err).Msg("failed get data by uID")
		return "", err
	}

	return data.(string), err
}

func (sr *cacheClient) DeleteDataByID(ctx context.Context, uID string) error {
	err := sr.RedisConnection.DeleteKeyValue(uID)
	if err != nil {
		return err
	}
	return nil
}
