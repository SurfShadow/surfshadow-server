package usecases

import (
	"fmt"

	"github.com/SurfShadow/surfshadow-server/internal/application/dto"
	"github.com/SurfShadow/surfshadow-server/internal/application/repositories"
	"github.com/SurfShadow/surfshadow-server/internal/domain/entities"
	"github.com/SurfShadow/surfshadow-server/pkg/logger"
)

const (
	errInvalidInputAllFieldsRequired = "invalid input: all fields are required"
	errInvalidInputIDRequired        = "invalid input: ID is required"
	errFailedToCreateProxyClient     = "failed to create proxy client"
	errFailedToGetAllProxyClients    = "failed to get all proxy clients"
	errFailedToGetProxyClientByID    = "failed to get proxy client by ID"
	errFailedToUpdateProxyClient     = "failed to update proxy client"
	errFailedDeleteProxyClient       = "failed to delete proxy client"
)

type ProxyClientUseCase struct {
	repo repositories.ProxyClientRepository
}

func NewProxyClientUseCase(repo repositories.ProxyClientRepository) *ProxyClientUseCase {
	logger.Instance.Debug("initializing ProxyClientUseCase")
	return &ProxyClientUseCase{
		repo: repo,
	}
}

func (uc *ProxyClientUseCase) CreateProxyClient(clientDTO *dto.ProxyClientDTO) (*dto.ProxyClientDTO, error) {
	logger.Instance.Infof("creating proxy client by proxy client DTO: %+v", clientDTO)

	if clientDTO.Title == "" || clientDTO.OS == "" || clientDTO.DownloadLink == "" {
		logger.Instance.Error(errInvalidInputAllFieldsRequired)
		return nil, fmt.Errorf(errInvalidInputAllFieldsRequired)
	}

	client := &entities.ProxyClient{
		Title:        clientDTO.Title,
		OS:           clientDTO.OS,
		DownloadLink: clientDTO.DownloadLink,
	}

	createdClient, err := uc.repo.Create(client)
	if err != nil {
		logger.Instance.Errorf("%s: %v", errFailedToCreateProxyClient, err)
		return nil, fmt.Errorf("%s: %w", errFailedToCreateProxyClient, err)
	}

	logger.Instance.Debugf("created client: %+v", createdClient)

	return &dto.ProxyClientDTO{
		ID:           createdClient.ID,
		Title:        createdClient.Title,
		OS:           createdClient.OS,
		DownloadLink: createdClient.DownloadLink,
	}, nil
}

func (uc *ProxyClientUseCase) GetAllProxyClients() ([]*dto.ProxyClientDTO, error) {
	logger.Instance.Info("getting all proxy clients")

	clients, err := uc.repo.GetAll()
	if err != nil {
		logger.Instance.Errorf("%s: %v", errFailedToGetAllProxyClients, err)
		return nil, fmt.Errorf("%s: %w", errFailedToGetAllProxyClients, err)
	}

	for _, client := range clients {
		logger.Instance.Debugf("fetched client: %+v", client)
	}

	var clientDTOs = make([]*dto.ProxyClientDTO, 0)
	for _, client := range clients {
		clientDTOs = append(clientDTOs, &dto.ProxyClientDTO{
			ID:           client.ID,
			Title:        client.Title,
			OS:           client.OS,
			DownloadLink: client.DownloadLink,
			CreatedAt:    client.CreatedAt,
			UpdatedAt:    client.UpdatedAt,
		})
	}

	logger.Instance.Debugf("fetched %d proxy clients", len(clientDTOs))

	return clientDTOs, nil
}

func (uc *ProxyClientUseCase) GetProxyClientByID(id int64) (*dto.ProxyClientDTO, error) {
	logger.Instance.Infof("getting proxy client by ID: %d", id)

	client, err := uc.repo.GetByID(id)
	if err != nil {
		logger.Instance.Errorf("%s: %v", errFailedToGetProxyClientByID, err)
		return nil, fmt.Errorf("%s: %w", errFailedToGetProxyClientByID, err)
	}

	logger.Instance.Debugf("fetched client: %+v", client)

	return &dto.ProxyClientDTO{
		ID:           client.ID,
		Title:        client.Title,
		OS:           client.OS,
		DownloadLink: client.DownloadLink,
	}, nil
}

func (uc *ProxyClientUseCase) UpdateProxyClient(clientDTO *dto.ProxyClientDTO) error {
	logger.Instance.Infof("updating client by proxy client DTO: %+v", clientDTO)

	if clientDTO.ID == 0 {
		logger.Instance.Error(errInvalidInputIDRequired)
		return fmt.Errorf(errInvalidInputIDRequired)
	}

	client := &entities.ProxyClient{
		ID:           clientDTO.ID,
		Title:        clientDTO.Title,
		OS:           clientDTO.OS,
		DownloadLink: clientDTO.DownloadLink,
	}

	logger.Instance.Debugf("updating client: %+v", client)

	err := uc.repo.Update(client)
	if err != nil {
		logger.Instance.Errorf("%s: %v", errFailedToUpdateProxyClient, err)
		return fmt.Errorf("%s: %v", errFailedToUpdateProxyClient, err)
	}

	logger.Instance.Infof("successfully updated proxy client with ID: %d", clientDTO.ID)

	return nil
}

func (uc *ProxyClientUseCase) DeleteProxyClient(id int64) error {
	logger.Instance.Infof("deleting client by ID: %d", id)

	if id == 0 {
		logger.Instance.Error(errInvalidInputIDRequired)
		return fmt.Errorf(errInvalidInputIDRequired)
	}

	logger.Instance.Debugf("deleting client with ID: %d", id)

	err := uc.repo.Delete(id)
	if err != nil {
		logger.Instance.Errorf("%s: %v", errFailedDeleteProxyClient, err)
		return fmt.Errorf("%s: %v", errFailedDeleteProxyClient, err)
	}

	logger.Instance.Infof("successfully deleted proxy client with ID: %d", id)

	return nil
}
