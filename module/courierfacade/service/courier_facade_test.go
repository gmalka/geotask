package service

import (
	"context"
	"reflect"
	"testing"

	"gitlab.com/ptflp/geotask/mocks"
	"gitlab.com/ptflp/geotask/module/courier/models"
	orderModels "gitlab.com/ptflp/geotask/module/order/models"
	cservice "gitlab.com/ptflp/geotask/module/courier/service"
	cfm "gitlab.com/ptflp/geotask/module/courierfacade/models"
	oservice "gitlab.com/ptflp/geotask/module/order/service"
)

func TestCourierFacade_MoveCourier(t *testing.T) {
	var (
		lat, lng, radius float64
		direction, zoom int
	)
	courierService := mocks.NewCourierer(t)
	orderService := mocks.NewOrderer(t)

	direction = 1
	zoom = 13
	lat = 10.0
	lng = 11.0
	radius = 5.0 * float64(19.0 - zoom)
	point := models.Point{Lat: lat, Lng: lng}
	courier := models.Courier{Score: 0, Location: point}
	wantedCourier := models.Courier{Score: 2, Location: point}

	courierService.On("GetCourier", context.Background()).Return(&courier, nil)
	orderService.On("DeleteByRadius", context.Background(), lat, lng, radius, "m").Return(2, nil)
	courierService.On("MoveCourier", wantedCourier, direction, zoom).Return(nil)

	type fields struct {
		courierService cservice.Courierer
		orderService   oservice.Orderer
	}
	type args struct {
		ctx       context.Context
		direction int
		zoom      int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Regular Move Courier",
			fields: fields{
				courierService: courierService,
				orderService: orderService,
			},
			args: args{
				ctx: context.Background(),
				direction: direction,
				zoom: zoom,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := CourierFacade{
				courierService: tt.fields.courierService,
				orderService:   tt.fields.orderService,
			}
			c.MoveCourier(tt.args.ctx, tt.args.direction, tt.args.zoom)
		})
	}
}

func TestCourierFacade_GetStatus(t *testing.T) {
	var (
		lat, lng float64
	)
	courierService := mocks.NewCourierer(t)
	orderService := mocks.NewOrderer(t)

	lat = 10.0
	lng = 11.0

	point := models.Point{Lat: lat, Lng: lng}
	courier := &models.Courier{Score: 0, Location: point}
	order := []orderModels.Order{{}}

	courierService.On("GetCourier", context.Background()).Return(courier, nil)
	orderService.On("GetByRadius", context.Background(), lng, lat, float64(CourierVisibilityRadius), "m").Return(order, nil)

	returnedStatus := cfm.CourierStatus{Courier: *courier, Orders: order}

	type fields struct {
		courierService cservice.Courierer
		orderService   oservice.Orderer
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   cfm.CourierStatus
	}{
		{
			name: "Regular GetStatus",
			fields: fields{
				courierService: courierService,
				orderService: orderService,
			},
			args: args{context.Background()},
			want: returnedStatus,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := CourierFacade{
				courierService: tt.fields.courierService,
				orderService:   tt.fields.orderService,
			}
			if got := c.GetStatus(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CourierFacade.GetStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}