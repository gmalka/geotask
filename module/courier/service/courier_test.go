package service

import (
	"context"
	"errors"
	"math"
	"reflect"
	"testing"

	"gitlab.com/ptflp/geotask/geo"
	"gitlab.com/ptflp/geotask/mocks"
	"gitlab.com/ptflp/geotask/module/courier/models"
	"gitlab.com/ptflp/geotask/module/courier/storage"
)

func TestNewCourierService(t *testing.T) {
	courierStorage := mocks.NewCourierStorager(t)

	point := models.Point{Lat: 9.00, Lng: 10.00}
	courier := models.Courier{Score: 0, Location: point}
	courierStorage.On("Save", context.Background(), courier).Return(nil)

	allowedZone := mocks.NewPolygonChecker(t)
	allowedZone.On("RandomPoint").Return(geo.Point(point))

	disbledZonesFirst := mocks.NewPolygonChecker(t)
	disbledZonesFirst.On("Contains", geo.Point(point)).Return(false)
	disbledZonesSecond := mocks.NewPolygonChecker(t)
	disbledZonesSecond.On("Contains", geo.Point(point)).Return(false)
	disabledZone := make([]geo.PolygonChecker, 0, 2)
	disabledZone = append(disabledZone, disbledZonesFirst)
	disabledZone = append(disabledZone, disbledZonesSecond)

	type args struct {
		courierStorage storage.CourierStorager
		allowedZone    geo.PolygonChecker
		disbledZones   []geo.PolygonChecker
	}
	tests := []struct {
		name string
		args args
		want Courierer
	}{
		{
			name: "Standart Curier Service Create",
			args: args{
				courierStorage: courierStorage,
				allowedZone: allowedZone,
				disbledZones: disabledZone,
			},
			want: &CourierService{courierStorage: courierStorage, allowedZone: allowedZone, disabledZones: disabledZone},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCourierService(tt.args.courierStorage, tt.args.allowedZone, tt.args.disbledZones); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCourierService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCourierService_GetCourier(t *testing.T) {

	// first
	courierStorage := mocks.NewCourierStorager(t)
	disbledZonesFirst := mocks.NewPolygonChecker(t)
	disbledZonesSecond := mocks.NewPolygonChecker(t)
	allowedZone := mocks.NewPolygonChecker(t)

	point := models.Point{Lat: 9.00, Lng: 10.00}
	courier := &models.Courier{Score: 0, Location: point}
	
	courierStorage.On("GetOne", context.Background()).Return(courier, nil)
	allowedZone.On("Contains", geo.Point(point)).Return(true)
	disbledZonesFirst.On("Contains", geo.Point(point)).Return(false)
	disbledZonesSecond.On("Contains", geo.Point(point)).Return(false)
	disabledZone := make([]geo.PolygonChecker, 0, 2)
	disabledZone = append(disabledZone, disbledZonesFirst)
	disabledZone = append(disabledZone, disbledZonesSecond)

	// second
	courierStorageSecond := mocks.NewCourierStorager(t)

	point2 := models.Point{Lat: 10.00, Lng: 11.00}
	courier2 := &models.Courier{Score: 0, Location: point2}
	
	courierStorageSecond.On("GetOne", context.Background()).Return(courier2, nil)
	allowedZone.On("Contains", geo.Point(point2)).Return(true)
	allowedZone.On("RandomPoint").Return(geo.Point(point))
	disbledZonesFirst.On("Contains", geo.Point(point2)).Return(true)
	disbledZonesFirst.On("Contains", geo.Point(point)).Return(false)
	disbledZonesSecond.On("Contains", geo.Point(point)).Return(false)
	disabledZoneSecond := make([]geo.PolygonChecker, 0, 2)
	disabledZoneSecond = append(disabledZoneSecond, disbledZonesFirst)
	disabledZoneSecond = append(disabledZoneSecond, disbledZonesSecond)

	// third
	courierStorageThird := mocks.NewCourierStorager(t)

	courierStorageThird.On("GetOne", context.Background()).Return(nil, errors.New("Some error"))

	type fields struct {
		courierStorage storage.CourierStorager
		allowedZone    geo.PolygonChecker
		disabledZones  []geo.PolygonChecker
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.Courier
		wantErr bool
	}{
		{
			name: "Standart get",
			fields: fields {
				courierStorage: courierStorage,
				allowedZone: allowedZone,
				disabledZones: disabledZone,
			},
			args: args{context.Background()},
			want: courier,
			wantErr: false,
		},
		{
			name: "Not in allowed area get",
			fields: fields {
				courierStorage: courierStorageSecond,
				allowedZone: allowedZone,
				disabledZones: disabledZoneSecond,
			},
			args: args{context.Background()},
			want: courier2,
			wantErr: false,
		},
		{
			name: "Error get",
			fields: fields {
				courierStorage: courierStorageThird,
				allowedZone: allowedZone,
				disabledZones: disabledZoneSecond,
			},
			args: args{context.Background()},
			want: nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CourierService{
				courierStorage: tt.fields.courierStorage,
				allowedZone:    tt.fields.allowedZone,
				disabledZones:  tt.fields.disabledZones,
			}
			got, err := c.GetCourier(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("CourierService.GetCourier() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CourierService.GetCourier() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCourierService_MoveCourier(t *testing.T) {
	// first
	point := models.Point{Lat: 10.00, Lng: 11.00}
	courier := &models.Courier{Score: 0, Location: point}
	wantPoint := models.Point{Lat: 10.00 - (float64(0.001) / math.Pow(2, float64(10 - 14))), Lng: 11.00}
	wantCourier := &models.Courier{Score: 0, Location: wantPoint}
	
	courierStorage := mocks.NewCourierStorager(t)
	allowedZone := mocks.NewPolygonChecker(t)
	disbledZonesFirst := mocks.NewPolygonChecker(t)
	disbledZonesSecond := mocks.NewPolygonChecker(t)

	allowedZone.On("Contains", geo.Point(wantPoint)).Return(true)
	disbledZonesFirst.On("Contains", geo.Point(wantPoint)).Return(false)
	disbledZonesSecond.On("Contains", geo.Point(wantPoint)).Return(false)
	courierStorage.On("Save", context.Background(), *wantCourier).Return(nil)

	disabledZone := make([]geo.PolygonChecker, 0, 2)
	disabledZone = append(disabledZone, disbledZonesFirst)
	disabledZone = append(disabledZone, disbledZonesSecond)

	// second
	point2 := models.Point{Lat: 10.00, Lng: 11.00}
	courier2 := &models.Courier{Score: 0, Location: point2}
	accuracy := float64(0.001) / math.Pow(2, float64(10 - 14))
	if accuracy < 0 {
		accuracy *= -1
	}
	wantPoint2 := models.Point{Lat: 10.00 + accuracy, Lng: 11.00}

	courierStorage2 := mocks.NewCourierStorager(t)
	allowedZone2 := mocks.NewPolygonChecker(t)
	disbledZonesFirst2 := mocks.NewPolygonChecker(t)
	disbledZonesSecond2 := mocks.NewPolygonChecker(t)

	allowedZone2.On("Contains", geo.Point(wantPoint2)).Return(false)
	allowedZone2.On("RandomPoint").Return(geo.Point(point2))
	disbledZonesFirst2.On("Contains", geo.Point(point2)).Return(false)
	disbledZonesSecond2.On("Contains", geo.Point(point2)).Return(false)
	courierStorage2.On("Save", context.Background(), *courier2).Return(nil)

	disabledZone2 := make([]geo.PolygonChecker, 0, 2)
	disabledZone2 = append(disabledZone2, disbledZonesFirst2)
	disabledZone2 = append(disabledZone2, disbledZonesSecond2)

	type fields struct {
		courierStorage storage.CourierStorager
		allowedZone    geo.PolygonChecker
		disabledZones  []geo.PolygonChecker
	}
	type args struct {
		courier   models.Courier
		direction int
		zoom      int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Standart Move",
			fields: fields{
				courierStorage: courierStorage,
				allowedZone: allowedZone,
				disabledZones: disabledZone,
			},
			args: args{
				courier: *courier,
				direction: 1,
				zoom: 10,
			},
			wantErr: false,
		},
		{
			name: "Not contain Move",
			fields: fields{
				courierStorage: courierStorage2,
				allowedZone: allowedZone2,
				disabledZones: disabledZone2,
			},
			args: args{
				courier: *courier2,
				direction: 0,
				zoom: 10,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CourierService{
				courierStorage: tt.fields.courierStorage,
				allowedZone:    tt.fields.allowedZone,
				disabledZones:  tt.fields.disabledZones,
			}
			if err := c.MoveCourier(tt.args.courier, tt.args.direction, tt.args.zoom); (err != nil) != tt.wantErr {
				t.Errorf("CourierService.MoveCourier() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCheckPointIsAllowed(t *testing.T) {
	// first
	allowedZone := mocks.NewPolygonChecker(t)
	disbledZonesFirst := mocks.NewPolygonChecker(t)
	disbledZonesSecond := mocks.NewPolygonChecker(t)

	point := models.Point{Lat: 10.00, Lng: 11.00}

	allowedZone.On("Contains", geo.Point(point)).Return(true)
	disbledZonesFirst.On("Contains", geo.Point(point)).Return(false)
	disbledZonesSecond.On("Contains", geo.Point(point)).Return(false)

	disabledZone := make([]geo.PolygonChecker, 0, 2)
	disabledZone = append(disabledZone, disbledZonesFirst)
	disabledZone = append(disabledZone, disbledZonesSecond)

	// second
	disbledZonesSecond2 := mocks.NewPolygonChecker(t)
	disbledZonesSecond2.On("Contains", geo.Point(point)).Return(true)

	disabledZone2 := make([]geo.PolygonChecker, 0, 2)
	disabledZone2 = append(disabledZone2, disbledZonesFirst)
	disabledZone2 = append(disabledZone2, disbledZonesSecond2)

	type args struct {
		point         geo.Point
		allowedZone   geo.PolygonChecker
		disabledZones []geo.PolygonChecker
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Regular allowed check",
			args: args{
				point: geo.Point(point),
				allowedZone: allowedZone,
				disabledZones: disabledZone,
			},
			want: true,
		},
		{
			name: "False allowed check",
			args: args{
				point: geo.Point(point),
				allowedZone: allowedZone,
				disabledZones: disabledZone2,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckPointIsAllowed(tt.args.point, tt.args.allowedZone, tt.args.disabledZones); got != tt.want {
				t.Errorf("CheckPointIsAllowed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetRandomAllowedLocation(t *testing.T) {
	allowedZone := mocks.NewPolygonChecker(t)
	disbledZonesFirst := mocks.NewPolygonChecker(t)
	disbledZonesSecond := mocks.NewPolygonChecker(t)

	point := models.Point{Lat: 10.00, Lng: 11.00}

	allowedZone.On("RandomPoint").Return(geo.Point(point))
	disbledZonesFirst.On("Contains", geo.Point(point)).Return(false)
	disbledZonesSecond.On("Contains", geo.Point(point)).Return(false)

	disabledZone := make([]geo.PolygonChecker, 0, 2)
	disabledZone = append(disabledZone, disbledZonesFirst)
	disabledZone = append(disabledZone, disbledZonesSecond)

	type args struct {
		allowedZone   geo.PolygonChecker
		disabledZones []geo.PolygonChecker
	}
	tests := []struct {
		name string
		args args
		want geo.Point
	}{
		{
			name: "Regular get random allowed lcoation",
			args: args{
				allowedZone: allowedZone,
				disabledZones: disabledZone,
			},
			want: geo.Point(point),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetRandomAllowedLocation(tt.args.allowedZone, tt.args.disabledZones); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRandomAllowedLocation() = %v, want %v", got, tt.want)
			}
		})
	}
}
