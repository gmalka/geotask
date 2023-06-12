package service

import (
	"context"
	"log"
	"math"

	"gitlab.com/ptflp/geotask/geo"
	"gitlab.com/ptflp/geotask/module/courier/models"
	"gitlab.com/ptflp/geotask/module/courier/storage"
)

// Направления движения курьера
const (
	DirectionUp    = 0
	DirectionDown  = 1
	DirectionLeft  = 2
	DirectionRight = 3
)

const (
	DefaultCourierLat = 59.9311
	DefaultCourierLng = 30.3609
)

//go:generate mockery --name Courierer
type Courierer interface {
	GetCourier(ctx context.Context) (*models.Courier, error)
	MoveCourier(courier models.Courier, direction, zoom int) error
}

type CourierService struct {
	courierStorage storage.CourierStorager
	allowedZone    geo.PolygonChecker
	disabledZones  []geo.PolygonChecker
}

func NewCourierService(courierStorage storage.CourierStorager, allowedZone geo.PolygonChecker, disbledZones []geo.PolygonChecker) Courierer {
	point := geo.GetRandomAllowedLocation(allowedZone, disbledZones)
	err := courierStorage.Save(context.Background(), models.Courier{Score: 0, Location: models.Point{Lat: point.Lat, Lng: point.Lng}})
	if err != nil {
		log.Fatalf("NewCourierService: %s\n", err)
	}
	return &CourierService{courierStorage: courierStorage, allowedZone: allowedZone, disabledZones: disbledZones}
}

func (c *CourierService) GetCourier(ctx context.Context) (*models.Courier, error) {
	// получаем курьера из хранилища используя метод GetOne из storage/courier.go
	courier, err := c.courierStorage.GetOne(ctx)
	if err != nil {
		return nil, err
	}

	if !geo.CheckPointIsAllowed(geo.Point{Lat: courier.Location.Lat, Lng: courier.Location.Lng}, c.allowedZone, c.disabledZones) {
		point := geo.GetRandomAllowedLocation(c.allowedZone, c.disabledZones)
		courier.Location.Lat = point.Lat
		courier.Location.Lng = point.Lng
	}
	// проверяем, что курьер находится в разрешенной зоне
	// если нет, то перемещаем его в случайную точку в разрешенной зоне
	// сохраняем новые координаты курьера

	return courier, nil
}

// MoveCourier : direction - направление движения курьера, zoom - зум карты
func (c *CourierService) MoveCourier(courier models.Courier, direction, zoom int) error {
	// точность перемещения зависит от зума карты использовать формулу 0.001 / 2^(zoom - 14)
	// 14 - это максимальный зум карты
	accuracy := float64(0.001) / math.Pow(2, float64(zoom - 14))

	if accuracy < 0 {
		accuracy *= -1
	}

	switch direction {
	case 0:
		courier.Location.Lat = courier.Location.Lat + accuracy
	case 1:
		courier.Location.Lat = courier.Location.Lat - accuracy
	case 2:
		courier.Location.Lng = courier.Location.Lng - accuracy
	case 3:
		courier.Location.Lng = courier.Location.Lng + accuracy
	}

	// далее нужно проверить, что курьер не вышел за границы зоны
	// если вышел, то нужно переместить его в случайную точку внутри зоны
	if !geo.CheckPointIsAllowed(geo.Point{Lat: courier.Location.Lat, Lng: courier.Location.Lng}, c.allowedZone, c.disabledZones) {
		point := geo.GetRandomAllowedLocation(c.allowedZone, c.disabledZones)
		courier.Location.Lat = point.Lat
		courier.Location.Lng = point.Lng
	}

	// далее сохранить изменения в хранилище
	return c.courierStorage.Save(context.Background(), courier)
}
