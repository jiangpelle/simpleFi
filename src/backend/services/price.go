package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/streadway/amqp"
)

type PriceService struct {
	redisClient *redis.Client
	rabbitMQ    *amqp.Connection
}

type PriceUpdate struct {
	Token     string  `json:"token"`
	Price     float64 `json:"price"`
	Timestamp int64   `json:"timestamp"`
}

func NewPriceService(redisClient *redis.Client, rabbitMQ *amqp.Connection) *PriceService {
	return &PriceService{
		redisClient: redisClient,
		rabbitMQ:    rabbitMQ,
	}
}

func (s *PriceService) StartPriceUpdates() error {
	ch, err := s.rabbitMQ.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"price_updates",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %v", err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %v", err)
	}

	go func() {
		for msg := range msgs {
			var update PriceUpdate
			if err := json.Unmarshal(msg.Body, &update); err != nil {
				log.Printf("Error decoding price update: %v", err)
				continue
			}

			if err := s.UpdatePrice(update); err != nil {
				log.Printf("Error updating price: %v", err)
			}
		}
	}()

	return nil
}

func (s *PriceService) UpdatePrice(update PriceUpdate) error {
	ctx := context.Background()
	key := fmt.Sprintf("price:%s", update.Token)

	// Store price in Redis with expiration
	err := s.redisClient.Set(ctx, key, update.Price, 24*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("failed to set price in Redis: %v", err)
	}

	// Store price history
	historyKey := fmt.Sprintf("price_history:%s", update.Token)
	historyData := map[string]interface{}{
		"price":     update.Price,
		"timestamp": update.Timestamp,
	}

	err = s.redisClient.ZAdd(ctx, historyKey, &redis.Z{
		Score:  float64(update.Timestamp),
		Member: historyData,
	}).Err()
	if err != nil {
		return fmt.Errorf("failed to add price to history: %v", err)
	}

	// Keep only last 1000 price points
	s.redisClient.ZRemRangeByRank(ctx, historyKey, 0, -1001)

	return nil
}

func (s *PriceService) GetCurrentPrice(token string) (float64, error) {
	ctx := context.Background()
	key := fmt.Sprintf("price:%s", token)

	price, err := s.redisClient.Get(ctx, key).Float64()
	if err != nil {
		return 0, fmt.Errorf("failed to get price: %v", err)
	}

	return price, nil
}

func (s *PriceService) GetPriceHistory(token string, start, end int64) ([]PriceUpdate, error) {
	ctx := context.Background()
	key := fmt.Sprintf("price_history:%s", token)

	results, err := s.redisClient.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min: fmt.Sprintf("%d", start),
		Max: fmt.Sprintf("%d", end),
	}).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get price history: %v", err)
	}

	var updates []PriceUpdate
	for _, result := range results {
		var update PriceUpdate
		if err := json.Unmarshal([]byte(result), &update); err != nil {
			continue
		}
		updates = append(updates, update)
	}

	return updates, nil
}
