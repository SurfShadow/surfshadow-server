package repositories

import "github.com/SurfShadow/surfshadow-server/internal/domain/entities"

type ProxyClientRepository interface {
	Create(client *entities.ProxyClient) (*entities.ProxyClient, error)
	GetAll() ([]*entities.ProxyClient, error)
	GetByID(id int64) (*entities.ProxyClient, error)
	Update(client *entities.ProxyClient) error
	Delete(id int64) error
}
