package service

import (
	"context"
	"encoding/json"

	"example.com/GoService/pkg/models"

	"github.com/nats-io/stan.go"
	"github.com/rs/zerolog"
)

const batchSize = 9

type storage interface {
	WriteToDB(ctx context.Context, order []models.Order) (err error)
	GetDataByID(ctx context.Context, uID string) (order models.Order, err error)
}

type Service struct {
	log      zerolog.Logger
	natsConn stan.Conn
	storage  storage
}

func (s *Service) GetDataByID(ctx context.Context, uID string) (order models.Order, err error) {
	return s.storage.GetDataByID(ctx, uID)
}

func (s *Service) ReceiveMessage(subject string) (err error) {
	var (
		ctx    = context.Background()
		order  models.Order
		orders = make([]models.Order, 0)
		cnt    = 0
	)
	if _, err = s.natsConn.Subscribe(subject, func(msg *stan.Msg) {
		if err = json.Unmarshal(msg.Data, &order); err != nil {
			s.log.Fatal().Err(err).Send()
		}
		orders = append(orders, order)
		cnt++
		s.log.Info().Msgf("Msg num - %d", cnt)
		if len(orders) == batchSize {
			if err = s.storage.WriteToDB(ctx, orders); err != nil {
				s.log.Fatal().Err(err).Send()
			}
		}
	}); err != nil {
		return
	}

	return
}

func New(log zerolog.Logger, natsConn stan.Conn, strg storage) *Service {
	return &Service{
		log:      log,
		natsConn: natsConn,
		storage:  strg,
	}
}
