package service

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/mock"
	"gitlab.com/ptflp/geotask/geo"
	"gitlab.com/ptflp/geotask/mocks"
	models "gitlab.com/ptflp/geotask/module/order/models"
	"gitlab.com/ptflp/geotask/module/order/storage"
)

func TestOrderService_GetByRadius(t *testing.T) {
	orderStorage := mocks.NewOrderStorager(t)
	orderStorage.On("GetByRadius", context.Background(), 10.0, 11.0, 200.0, "m").Return([]models.Order{{}, {}}, nil)

	allowedZone := mocks.NewPolygonChecker(t)

	disbledZonesFirst := mocks.NewPolygonChecker(t)
	disbledZonesSecond := mocks.NewPolygonChecker(t)
	disabledZone := make([]geo.PolygonChecker, 0, 2)
	disabledZone = append(disabledZone, disbledZonesFirst)
	disabledZone = append(disabledZone, disbledZonesSecond)

	type fields struct {
		storage       storage.OrderStorager
		allowedZone   geo.PolygonChecker
		disabledZones []geo.PolygonChecker
	}
	type args struct {
		ctx    context.Context
		lng    float64
		lat    float64
		radius float64
		unit   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []models.Order
		wantErr bool
	}{
		{
			name: "Regular Get By Radius",
			fields: fields{
				storage:       orderStorage,
				allowedZone:   allowedZone,
				disabledZones: disabledZone,
			},
			args: args{
				ctx:    context.Background(),
				lng:    10.0,
				lat:    11.0,
				radius: 200.0,
				unit:   "m",
			},
			want:    []models.Order{{}, {}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &OrderService{
				storage:       tt.fields.storage,
				allowedZone:   tt.fields.allowedZone,
				disabledZones: tt.fields.disabledZones,
			}
			got, err := o.GetByRadius(tt.args.ctx, tt.args.lng, tt.args.lat, tt.args.radius, tt.args.unit)
			if (err != nil) != tt.wantErr {
				t.Errorf("OrderService.GetByRadius() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OrderService.GetByRadius() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderService_DeleteByRadius(t *testing.T) {
	orderStorage := mocks.NewOrderStorager(t)
	orderStorage.On("DeleteByRadius", context.Background(), 10.0, 11.0, 200.0, "m").Return(2, nil)

	allowedZone := mocks.NewPolygonChecker(t)

	disbledZonesFirst := mocks.NewPolygonChecker(t)
	disbledZonesSecond := mocks.NewPolygonChecker(t)
	disabledZone := make([]geo.PolygonChecker, 0, 2)
	disabledZone = append(disabledZone, disbledZonesFirst)
	disabledZone = append(disabledZone, disbledZonesSecond)

	type fields struct {
		storage       storage.OrderStorager
		allowedZone   geo.PolygonChecker
		disabledZones []geo.PolygonChecker
	}
	type args struct {
		ctx    context.Context
		lat    float64
		lng    float64
		radius float64
		unit   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "Regular Delete By Radius",
			fields: fields{
				storage:       orderStorage,
				allowedZone:   allowedZone,
				disabledZones: disabledZone,
			},
			args: args{
				ctx:    context.Background(),
				lng:    11.0,
				lat:    10.0,
				radius: 200.0,
				unit:   "m",
			},
			want:    2,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &OrderService{
				storage:       tt.fields.storage,
				allowedZone:   tt.fields.allowedZone,
				disabledZones: tt.fields.disabledZones,
			}
			got, err := o.DeleteByRadius(tt.args.ctx, tt.args.lat, tt.args.lng, tt.args.radius, tt.args.unit)
			if (err != nil) != tt.wantErr {
				t.Errorf("OrderService.DeleteByRadius() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("OrderService.DeleteByRadius() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderService_Save(t *testing.T) {
	orderStorage := mocks.NewOrderStorager(t)
	orderStorage.On("Save", context.Background(), models.Order{}, orderMaxAge).Return(nil)

	allowedZone := mocks.NewPolygonChecker(t)

	disbledZonesFirst := mocks.NewPolygonChecker(t)
	disbledZonesSecond := mocks.NewPolygonChecker(t)
	disabledZone := make([]geo.PolygonChecker, 0, 2)
	disabledZone = append(disabledZone, disbledZonesFirst)
	disabledZone = append(disabledZone, disbledZonesSecond)

	type fields struct {
		storage       storage.OrderStorager
		allowedZone   geo.PolygonChecker
		disabledZones []geo.PolygonChecker
	}
	type args struct {
		ctx   context.Context
		order models.Order
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Regular Save",
			fields: fields{
				storage:       orderStorage,
				allowedZone:   allowedZone,
				disabledZones: disabledZone,
			},
			args: args{
				ctx:   context.Background(),
				order: models.Order{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &OrderService{
				storage:       tt.fields.storage,
				allowedZone:   tt.fields.allowedZone,
				disabledZones: tt.fields.disabledZones,
			}
			if err := o.Save(tt.args.ctx, tt.args.order); (err != nil) != tt.wantErr {
				t.Errorf("OrderService.Save() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOrderService_GetCount(t *testing.T) {
	orderStorage := mocks.NewOrderStorager(t)
	orderStorage.On("GetCount", context.Background()).Return(12, nil)

	allowedZone := mocks.NewPolygonChecker(t)

	disbledZonesFirst := mocks.NewPolygonChecker(t)
	disbledZonesSecond := mocks.NewPolygonChecker(t)
	disabledZone := make([]geo.PolygonChecker, 0, 2)
	disabledZone = append(disabledZone, disbledZonesFirst)
	disabledZone = append(disabledZone, disbledZonesSecond)

	type fields struct {
		storage       storage.OrderStorager
		allowedZone   geo.PolygonChecker
		disabledZones []geo.PolygonChecker
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "Regular Get Count",
			fields: fields{
				storage:       orderStorage,
				allowedZone:   allowedZone,
				disabledZones: disabledZone,
			},
			args: args{
				ctx: context.Background(),
			},
			want:    12,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &OrderService{
				storage:       tt.fields.storage,
				allowedZone:   tt.fields.allowedZone,
				disabledZones: tt.fields.disabledZones,
			}
			got, err := o.GetCount(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("OrderService.GetCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("OrderService.GetCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderService_RemoveOldOrders(t *testing.T) {
	orderStorage := mocks.NewOrderStorager(t)
	orderStorage.On("RemoveOldOrders", context.Background(), orderMaxAge).Return(nil)

	allowedZone := mocks.NewPolygonChecker(t)

	disbledZonesFirst := mocks.NewPolygonChecker(t)
	disbledZonesSecond := mocks.NewPolygonChecker(t)
	disabledZone := make([]geo.PolygonChecker, 0, 2)
	disabledZone = append(disabledZone, disbledZonesFirst)
	disabledZone = append(disabledZone, disbledZonesSecond)

	type fields struct {
		storage       storage.OrderStorager
		allowedZone   geo.PolygonChecker
		disabledZones []geo.PolygonChecker
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Regular Remove Old Orders",
			fields: fields{
				storage: orderStorage,
				allowedZone: allowedZone,
				disabledZones: disabledZone,
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &OrderService{
				storage:       tt.fields.storage,
				allowedZone:   tt.fields.allowedZone,
				disabledZones: tt.fields.disabledZones,
			}
			if err := o.RemoveOldOrders(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("OrderService.RemoveOldOrders() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOrderService_GenerateOrder(t *testing.T) {
	orderStorage := mocks.NewOrderStorager(t)

	allowedZone := mocks.NewPolygonChecker(t)
	point := geo.Point{Lat: 10.0, Lng: 11.0}
	allowedZone.On("RandomPoint").Return(point)

	disbledZonesFirst := mocks.NewPolygonChecker(t)
	disbledZonesSecond := mocks.NewPolygonChecker(t)
	disbledZonesFirst.On("Contains", point).Return(false)
	disbledZonesSecond.On("Contains", point).Return(false)
	disabledZone := make([]geo.PolygonChecker, 0, 2)
	disabledZone = append(disabledZone, disbledZonesFirst)
	disabledZone = append(disabledZone, disbledZonesSecond)

	orderStorage.On("Save", context.Background(), mock.AnythingOfType("models.Order"), orderMaxAge).Return(nil)

	type fields struct {
		storage       storage.OrderStorager
		allowedZone   geo.PolygonChecker
		disabledZones []geo.PolygonChecker
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Regular Generate Order",
			fields: fields{
				storage: orderStorage,
				allowedZone: allowedZone,
				disabledZones: disabledZone,
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &OrderService{
				storage:       tt.fields.storage,
				allowedZone:   tt.fields.allowedZone,
				disabledZones: tt.fields.disabledZones,
			}
			if err := o.GenerateOrder(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("OrderService.GenerateOrder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
