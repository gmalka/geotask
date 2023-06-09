package service

import (
	"context"
	"math/rand"
	"time"

	"gitlab.com/ptflp/geotask/geo"
	"gitlab.com/ptflp/geotask/module/order/models"
	"gitlab.com/ptflp/geotask/module/order/storage"
)

const (
	minDeliveryPrice = 100.00
	maxDeliveryPrice = 500.00

	maxOrderPrice = 3000.00
	minOrderPrice = 1000.00

	orderMaxAge = 2 * time.Minute
)

var id int64

type Orderer interface {
	GetByRadius(ctx context.Context, lng, lat, radius float64, unit string) ([]models.Order, error) // возвращает заказы через метод storage.GetByRadius
	Save(ctx context.Context, order models.Order) error                                             // сохраняет заказ через метод storage.Save с заданным временем жизни OrderMaxAge
	GetCount(ctx context.Context) (int, error)                                                      // возвращает количество заказов через метод storage.GetCount
	RemoveOldOrders(ctx context.Context) error                                                      // удаляет старые заказы через метод storage.RemoveOldOrders с заданным временем жизни OrderMaxAge
	GenerateOrder(ctx context.Context) error                                                        // генерирует заказ в случайной точке из разрешенной зоны, с уникальным id, ценой и ценой доставки
}

// OrderService реализация интерфейса Orderer
// в нем должны быть методы GetByRadius, Save, GetCount, RemoveOldOrders, GenerateOrder
// данный сервис отвечает за работу с заказами
type OrderService struct {
	storage       storage.OrderStorager
	allowedZone   geo.PolygonChecker
	disabledZones []geo.PolygonChecker
}

func NewOrderService(storage storage.OrderStorager, allowedZone geo.PolygonChecker, disallowedZone []geo.PolygonChecker) Orderer {
	return &OrderService{storage: storage, allowedZone: allowedZone, disabledZones: disallowedZone}
}

func (o *OrderService) GetByRadius(ctx context.Context, lng, lat, radius float64, unit string) ([]models.Order, error) {
	return o.storage.GetByRadius(ctx, lng, lat, radius, unit)
}

func (o *OrderService) Save(ctx context.Context, order models.Order) error {
	return o.storage.Save(ctx, order, orderMaxAge)
}

func (o *OrderService) GetCount(ctx context.Context) (int, error) {
	return o.storage.GetCount(ctx)
}

func (o *OrderService) RemoveOldOrders(ctx context.Context) error {
	return o.storage.RemoveOldOrders(ctx, orderMaxAge)
}

func (o *OrderService) GenerateOrder(ctx context.Context) error {
	rand.Seed(time.Now().Unix())
	deliveryPrice := minDeliveryPrice + rand.Float64() * (maxDeliveryPrice - minDeliveryPrice)
	orderPrice := minDeliveryPrice + rand.Float64() * (maxOrderPrice - minDeliveryPrice)
	point := geo.GetRandomAllowedLocation(o.allowedZone, o.disabledZones)

	order := models.Order{ID: id, Price: orderPrice, DeliveryPrice: deliveryPrice, CreatedAt: time.Now(), Lat: point.Lat, Lng: point.Lng}
	id++

	return o.Save(ctx, order)
}