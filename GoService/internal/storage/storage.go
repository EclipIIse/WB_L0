package storage

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"example.com/GoService/pkg/models"

	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog"
)

type CacheClient interface {
	CachedData(ctx context.Context, data string, uID string) error
	GetDataByID(ctx context.Context, uID string) (string, error)
	DeleteDataByID(ctx context.Context, uID string) error
}

type Storage struct {
	log   zerolog.Logger
	conn  *pgx.Conn
	cache CacheClient
}

func (s *Storage) WriteToDB(ctx context.Context, orders []models.Order) (err error) {
	var (
		jsonOrder []byte
		query     = "INSERT INTO wb.orders (id, data) VALUES ($1, $2)"
	)

	batch := &pgx.Batch{}
	for _, order := range orders {
		uID := order.OrderUid
		jsonOrder, err = json.Marshal(order)
		if err != nil {
			s.log.Error().Err(err).Send()
			return
		}

		batch.Queue(query, uID, jsonOrder)
		if err = s.cache.CachedData(ctx, string(jsonOrder), uID); err != nil {
			s.log.Error().Err(err).Msg("failed add data to cache")
			return
		}
	}

	ctxDb, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	if err = s.conn.SendBatch(ctxDb, batch).Close(); err != nil {
		if errors.Is(err, context.Canceled) {
			s.log.Error().Err(err).Msg("query time out")
			return
		}
		s.log.Error().Err(err).Send()
		return
	}

	return
}

func (s *Storage) GetDataByID(ctx context.Context, uID string) (order models.Order, err error) {
	query := "SELECT data FROM wb.orders WHERE id=$1"

	rawData, err := s.cache.GetDataByID(ctx, uID)
	if err != nil {
		s.log.Error().Err(err).Send()
	}

	if len(rawData) != 0 {
		if err = json.Unmarshal([]byte(rawData), &order); err != nil {
			s.log.Error().Err(err).Send()
			return
		}
		s.log.Info().Msg("get data from cache")
		return
	}

	ctxDb, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	if err = s.conn.QueryRow(ctxDb, query, uID).Scan(&order); err != nil {
		if errors.Is(err, context.Canceled) {
			s.log.Error().Err(err).Msg("query time out")
			return
		}
		s.log.Error().Err(err).Send()
		return
	}

	jsonOrder, err := json.Marshal(order)
	if err != nil {
		s.log.Error().Err(err).Send()
	}

	if err = s.cache.CachedData(ctx, string(jsonOrder), uID); err != nil {
		s.log.Error().Err(err).Msg("failed add data to cache")
		return
	}

	s.log.Info().Msg("get data from database")

	return
}

func New(log zerolog.Logger, pgConn *pgx.Conn, cClient CacheClient) *Storage {
	return &Storage{
		log:   log,
		conn:  pgConn,
		cache: cClient,
	}
}
