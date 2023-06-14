package service

import (
	"context"
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

	courierStorage := mocks.NewCourierStorager(t)

	point := models.Point{Lat: 9.00, Lng: 10.00}
	courier := &models.Courier{Score: 0, Location: point}
	courierStorage.On("GetOne", context.Background()).Return(courier, nil)

	allowedZone := mocks.NewPolygonChecker(t)
	allowedZone.On("Contains", geo.Point(point)).Return(true)

	disbledZonesFirst := mocks.NewPolygonChecker(t)
	disbledZonesFirst.On("Contains", geo.Point(point)).Return(false)
	disbledZonesSecond := mocks.NewPolygonChecker(t)
	disbledZonesSecond.On("Contains", geo.Point(point)).Return(false)
	disabledZone := make([]geo.PolygonChecker, 0, 2)
	disabledZone = append(disabledZone, disbledZonesFirst)
	disabledZone = append(disabledZone, disbledZonesSecond)

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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
	type args struct {
		allowedZone   geo.PolygonChecker
		disabledZones []geo.PolygonChecker
	}
	tests := []struct {
		name string
		args args
		want geo.Point
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetRandomAllowedLocation(tt.args.allowedZone, tt.args.disabledZones); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRandomAllowedLocation() = %v, want %v", got, tt.want)
			}
		})
	}
}
