package repositories

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/SurfShadow/surfshadow-server/internal/application/repositories"
	"github.com/SurfShadow/surfshadow-server/internal/domain/entities"
	"github.com/SurfShadow/surfshadow-server/pkg/logger"
)

const (
	ErrFailedToInsertProxyClient    = "failed to insert proxy client"
	ErrFailedToFetchProxyClients    = "failed to fetch proxy clients"
	ErrFailedToFetchProxyClientByID = "failed to fetch proxy client by ID"
	ErrFailedToUpdateProxyClient    = "failed to update proxy client"
	ErrFailedToFetchRowsAffected    = "failed to fetch rows affected"
	ErrNoRowsAffected               = "no rows affected, proxy client not found"
	ErrFailedToDeleteProxyClient    = "failed to delete proxy client"
)

type ProxyClientRepositoryImpl struct {
	db *sqlx.DB
}

func NewProxyClientRepository(db *sqlx.DB) repositories.ProxyClientRepository {
	logger.Instance.Debug("Initializing ProxyClientRepository")
	return &ProxyClientRepositoryImpl{
		db: db,
	}
}

func (r *ProxyClientRepositoryImpl) Create(client *entities.ProxyClient) (*entities.ProxyClient, error) {
	logger.Instance.Debugf("Inserting new proxy client: %+v", client)

	query := `
		INSERT INTO proxy_clients (title, os, download_link, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	row := r.db.QueryRow(query, client.Title, client.OS, client.DownloadLink)
	if err := row.Scan(&client.ID, &client.CreatedAt, &client.UpdatedAt); err != nil {
		logger.Instance.Errorf("%s: %v", ErrFailedToInsertProxyClient, err)
		return nil, fmt.Errorf("%s: %w", ErrFailedToInsertProxyClient, err)
	}

	logger.Instance.Debugf("Inserted proxy client: %+v", client)

	return client, nil
}

func (r *ProxyClientRepositoryImpl) GetAll() ([]*entities.ProxyClient, error) {
	logger.Instance.Debug("Fetching all proxy clients")

	query := `
		SELECT id, title, os, download_link, created_at, updated_at
		FROM proxy_clients
	`

	var clients []*entities.ProxyClient
	if err := r.db.Select(&clients, query); err != nil {
		logger.Instance.Errorf("%s: %v", ErrFailedToFetchProxyClients, err)
		return nil, fmt.Errorf("%s: %w", ErrFailedToFetchProxyClients, err)
	}

	for _, client := range clients {
		logger.Instance.Debugf("Fetched client: %+v", client)
	}

	logger.Instance.Debugf("Fetched %d proxy clients", len(clients))

	return clients, nil
}

func (r *ProxyClientRepositoryImpl) GetByID(id int64) (*entities.ProxyClient, error) {
	logger.Instance.Debugf("Fetching proxy client by ID: %d", id)

	query := `
		SELECT id, title, os, download_link, created_at, updated_at
		FROM proxy_clients
		WHERE id = $1
	`

	var client entities.ProxyClient
	if err := r.db.Get(&client, query, id); err != nil {
		logger.Instance.Errorf("%s: %v", ErrFailedToFetchProxyClientByID, err)
		return nil, fmt.Errorf("%s: %w", ErrFailedToFetchProxyClientByID, err)
	}

	logger.Instance.Debugf("Fetched client: %+v", client)

	return &client, nil
}

func (r *ProxyClientRepositoryImpl) Update(client *entities.ProxyClient) error {
	logger.Instance.Debugf("Updating proxy client: %+v", client)

	query := `
		UPDATE proxy_clients
		SET title = $1, os = $2, download_link = $3, updated_at = NOW()
		WHERE id = $4
	`

	result, err := r.db.Exec(query, client.Title, client.OS, client.DownloadLink, client.ID)
	if err != nil {
		logger.Instance.Errorf("%s: %v", ErrFailedToUpdateProxyClient, err)
		return fmt.Errorf("%s: %w", ErrFailedToUpdateProxyClient, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Instance.Errorf("%s: %v", ErrFailedToFetchRowsAffected, err)
		return fmt.Errorf("%s: %w", ErrFailedToFetchRowsAffected, err)
	}

	if rowsAffected == 0 {
		logger.Instance.Errorf(ErrNoRowsAffected)
		return fmt.Errorf(ErrNoRowsAffected)
	}

	logger.Instance.Debugf("Updated proxy client: %+v", client)

	return nil
}

func (r *ProxyClientRepositoryImpl) Delete(id int64) error {
	logger.Instance.Infof("Deleting proxy client by ID: %d", id)

	query := `
		DELETE FROM proxy_clients
		WHERE id = $1
	`

	result, err := r.db.Exec(query, id)
	if err != nil {
		logger.Instance.Errorf("%s: %v", ErrFailedToDeleteProxyClient, err)
		return fmt.Errorf("%s: %w", ErrFailedToDeleteProxyClient, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Instance.Errorf("%s: %v", ErrFailedToFetchRowsAffected, err)
		return fmt.Errorf("%s: %w", ErrFailedToFetchRowsAffected, err)
	}

	if rowsAffected == 0 {
		logger.Instance.Error(ErrNoRowsAffected)
		return fmt.Errorf(ErrNoRowsAffected)
	}

	logger.Instance.Debugf("Deleted proxy client with ID: %d", id)

	return nil
}
