package storage

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis"
	"gitlab.com/ptflp/geotask/module/courier/models"
)

//go:generate mockery --name CourierStorager
type CourierStorager interface {
	Save(ctx context.Context, courier models.Courier) error // сохранить курьера по ключу courier
	GetOne(ctx context.Context) (*models.Courier, error)    // получить курьера по ключу courier
}

type CourierStorage struct {
	storage *redis.Client
}

const key = "courier"

func NewCourierStorage(storage *redis.Client) CourierStorager {
	return &CourierStorage{storage: storage}
}

func (c CourierStorage) Save(ctx context.Context, courier models.Courier) error {
	b, err := json.Marshal(courier)
	if err != nil {
		return err
	}

	_, err = c.storage.Set(key, b, 0).Result()

	return err
}

func (c CourierStorage)GetOne(ctx context.Context) (*models.Courier, error) {
	var courier *models.Courier

	b, err := c.storage.Get(key).Bytes()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &courier)
	if err != nil {
		return nil, err
	}

	return courier, nil
}