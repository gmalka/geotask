package docs

import "gitlab.com/ptflp/geotask/module/courierfacade/models"

// добавить документацию для роута /api/status

// swagger:route GET /api/status status emptyRequest
// Получение статуса сервиса.
// responses:
//  200: getStatusResponse

// swagger:parameters emptyRequest
type emptyRequest struct {}

// swagger:response getStatusResponse
type getStatusResponse struct {
	// Информация о курьере и карте
	//
	// in:body
	Body models.CourierStatus
}