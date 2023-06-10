package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"gitlab.com/ptflp/geotask/module/order/models"
)

type OrderStorager interface {
	Save(ctx context.Context, order models.Order, maxAge time.Duration) error                       // сохранить заказ с временем жизни
	GetByID(ctx context.Context, orderID int) (*models.Order, error)                                // получить заказ по id
	GenerateUniqueID(ctx context.Context) (int64, error)                                            // сгенерировать уникальный id
	GetByRadius(ctx context.Context, lng, lat, radius float64, unit string) ([]models.Order, error) // получить заказы в радиусе от точки
	GetCount(ctx context.Context) (int, error)                                                      // получить количество заказов
	RemoveOldOrders(ctx context.Context, maxAge time.Duration) error   								 // удалить старые заказы по истечению времени maxAge
	DeleteByRadius(ctx context.Context, lng, lat, radius float64, unit string) (int, error)
}

const (
	key = "order"
	geoKey = "geoorder"
)

type OrderStorage struct {
	storage *redis.Client
}

func NewOrderStorage(storage *redis.Client) OrderStorager {
	return &OrderStorage{storage: storage}
}

func (o *OrderStorage) Save(ctx context.Context, order models.Order, maxAge time.Duration) error {
	// save with geo redis
	
	return o.saveOrderWithGeo(ctx, order, maxAge)
}

func (o *OrderStorage) DeleteByRadius(ctx context.Context, lng, lat, radius float64, unit string) (int, error){
	var err error
	var data []byte
	var ordersLocation []redis.GeoLocation

	// используем метод getOrdersByRadius для получения ID заказов в радиусе
	ordersLocation, err = o.getOrdersByRadius(ctx, lng, lat, radius, unit)
	// обратите внимание, что в случае отсутствия заказов в радиусе
	// метод getOrdersByRadius должен вернуть nil, nil (при ошибке redis.Nil)
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	// проходим по списку ID заказов и получаем данные о заказе
	// получаем данные о заказе по ID из redis по ключу order:ID
	count := 0
	for _, v := range ordersLocation {
		count++
		_, err = o.storage.Del(fmt.Sprintf("%s:%s", key, v.Name)).Result()
		if err != nil {
			return 0, err
		}

		_, err = o.storage.ZRem(geoKey, v.Name).Result()
		if err != nil {
			return 0, err
		}

		data, err = o.storage.Get(fmt.Sprintf("%s:%s", key, v.Name)).Bytes()
		if err != nil {
			continue
		}

		order := models.Order{}
		err = json.Unmarshal(data, &order)
		if err != nil {
			return 0, err
		}

		_, err = o.storage.ZRem(geoKey, float64(order.CreatedAt.UnixMicro())).Result()
		if err != nil {
			return 0, err
		}
	}

	return count, nil
}

func (o *OrderStorage) RemoveOldOrders(ctx context.Context, maxAge time.Duration) error {
	// получить ID всех старых ордеров, которые нужно удалить
	// используя метод ZRangeByScore
	// старые ордеры это те, которые были созданы две минуты назад
	// и более
	/**
	&redis.ZRangeBy{
		Max: использовать вычисление времени с помощью maxAge,
		Min: "0",
	}
	*/
	t := strconv.Itoa(int(time.Now().UnixMicro() - maxAge.Microseconds()))

	_, err := o.storage.ZRangeByScore(key, redis.ZRangeBy{Min: "0", Max: t}).Result()
	if err == redis.Nil {
		return nil
	}
	if err != nil {
		return err
	}

	// Проверить количество старых ордеров
	// удалить старые ордеры из redis используя метод ZRemRangeByScore где ключ "orders" min "-inf" max "(время создания старого ордера)"
	// удалять ордера по ключу не нужно, они будут удалены автоматически по истечению времени жизни
	o.storage.ZRemRangeByScore(key, "0", t)

	return nil
}

// 1686405326813278638
//        120000000000
// 1686405446813281196

func (o *OrderStorage) GetByID(ctx context.Context, orderID int) (*models.Order, error) {
	var err error
	var data []byte
	var order models.Order
	// получаем ордер из redis по ключу order:ID
	data, err = o.storage.Get(fmt.Sprintf("%s:%d", key, orderID)).Bytes()

	// проверяем что ордер не найден исключение redis.Nil, в этом случае возвращаем nil, nil
	if err == redis.Nil {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	// десериализуем ордер из json
	err = json.Unmarshal(data, &order)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (o *OrderStorage) saveOrderWithGeo(ctx context.Context, order models.Order, maxAge time.Duration) error {
	var err error
	var data []byte

	// сериализуем ордер в json
	data, err = json.Marshal(order)
	if err != nil {
		return err
	}

	id, err := o.GenerateUniqueID(ctx)
	if err != nil {
		return err
	}
	// сохраняем ордер в json redis по ключу order:ID с временем жизни maxAge
	o.storage.Set(fmt.Sprintf("%s:%d", key, id), data, maxAge)

	// добавляем ордер в гео индекс используя метод GeoAdd где Name - это ключ ордера, а Longitude и Latitude - координаты
	o.storage.GeoAdd(geoKey, &redis.GeoLocation{Name: fmt.Sprintf("%d", id), Longitude: order.Lng, Latitude: order.Lat})

	// zset сохраняем ордер для получения количества заказов со сложностью O(1)
	o.storage.ZAdd(key, redis.Z{Score: float64(order.CreatedAt.UnixMicro()), Member: data})
	// Score - время создания ордера
	// float64(time.Now().UnixMicro())

	return nil
}

func (o *OrderStorage) GetCount(ctx context.Context) (int, error) {
	// получить количество ордеров в упорядоченном множестве используя метод ZCard
	count, err := o.storage.ZCard(key).Result()
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (o *OrderStorage) GetByRadius(ctx context.Context, lng, lat, radius float64, unit string) ([]models.Order, error) {
	var err error
	var orders []models.Order
	var data []byte
	var ordersLocation []redis.GeoLocation

	// используем метод getOrdersByRadius для получения ID заказов в радиусе
	ordersLocation, err = o.getOrdersByRadius(ctx, lng, lat, radius, unit)
	// обратите внимание, что в случае отсутствия заказов в радиусе
	// метод getOrdersByRadius должен вернуть nil, nil (при ошибке redis.Nil)
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	orders = make([]models.Order, 0, len(ordersLocation))
	// проходим по списку ID заказов и получаем данные о заказе
	// получаем данные о заказе по ID из redis по ключу order:ID
	for _, v := range ordersLocation {
		data, err = o.storage.Get(fmt.Sprintf("%s:%s", key, v.Name)).Bytes()
		if err != nil {
			continue
		}

		order := models.Order{}
		err = json.Unmarshal(data, &order)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func (o *OrderStorage) getOrdersByRadius(ctx context.Context, lng, lat, radius float64, unit string) ([]redis.GeoLocation, error) {
	// в данном методе мы получаем список ордеров в радиусе от точки
	// возвращаем список ордеров с координатами и расстоянием до точки
	/**
	&redis.GeoRadiusQuery{
		Radius:      radius,
		Unit:        unit,
		WithCoord:   true,
		WithDist:    true,
		WithGeoHash: true,
	}
	*/

	result, err := o.storage.GeoRadius(geoKey, lng, lat, &redis.GeoRadiusQuery {
		Radius:      radius,
		Unit:        unit,
		WithCoord:   true,
		WithDist:    true,
		WithGeoHash: true,
	}).Result()

	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, redis.Nil
	}
	
	return result, nil
}

func (o *OrderStorage) GenerateUniqueID(ctx context.Context) (int64, error) {
	var err error
	var id int64

	// генерируем уникальный ID для ордера
	// используем для этого redis incr по ключу order:id
	id, err = o.storage.Incr(fmt.Sprintf("%s:%s", key, "id")).Result()
	if err != nil {
		return 0, err
	}

	return id, nil
}