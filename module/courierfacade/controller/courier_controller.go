package controller

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/ptflp/geotask/module/courierfacade/service"
)

type CourierController struct {
	courierService service.CourierFacer
}

func NewCourierController(courierService service.CourierFacer) *CourierController {
	return &CourierController{courierService: courierService}
}

func (c *CourierController) GetStatus(ctx *gin.Context) {
	// установить задержку в 50 миллисекунд
	time.Sleep(time.Millisecond * 50)

	// получить статус курьера из сервиса courierService используя метод GetStatus
	status := c.courierService.GetStatus(ctx)
	// отправить статус курьера в ответ
	ctx.JSON(http.StatusOK, status)
}

func (c *CourierController) MoveCourier(m webSocketMessage) {
	var cm CourierMove
	var err error
	// получить данные из m.Data и десериализовать их в структуру CourierMove
	b, err := json.Marshal(m.Data)
	if err != nil {
		log.Println(err)
	}

	err  = json.Unmarshal(b, &cm)
	if err != nil {
		log.Println(err)
	}

	c.courierService.MoveCourier(context.Background(), cm.Direction, cm.Zoom)
	// вызвать метод MoveCourier у courierService
}