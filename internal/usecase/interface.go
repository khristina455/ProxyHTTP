package usecase

import (
	"Proxy/internal/models"
	"net/http"
)

type UsecaseI interface {
	SaveRequest(request *http.Request) (int, error)
	SaveResponse(requestId int, response *http.Response) (models.Response, error)
	AllRequests() ([]models.Request, error)
	GetRequest(id int) (models.Request, error)
	RepeatRequest(id int) (models.Response, error)
}
