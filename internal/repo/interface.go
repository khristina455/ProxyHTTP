package repo

import "Proxy/internal/models"

type RepoI interface {
	SaveRequest(request models.Request) (int, error)
	AllRequests() ([]models.Request, error)
	SaveResponse(requestId int, response models.Response) error
	GetRequest(id int) (models.Request, error)
}
