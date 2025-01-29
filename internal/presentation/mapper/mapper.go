package mapper

import (
	"time"

	appdto "github.com/SurfShadow/surfshadow-server/internal/application/dto"
	"github.com/SurfShadow/surfshadow-server/internal/presentation/dto/proxy_client"
	"github.com/SurfShadow/surfshadow-server/pkg/logger"
)

func MapRequestToAppDTO(req proxy_client.ProxyClientRequest) appdto.ProxyClientDTO {
	return appdto.ProxyClientDTO{
		Title:        req.Title,
		OS:           req.OS,
		DownloadLink: req.DownloadLink,
	}
}

func MapAppDTOToResponse(aDTO appdto.ProxyClientDTO) proxy_client.ProxyClientResponse {
	logger.Instance.Debugf("Mapping app dto to response: %+v", aDTO.CreatedAt)

	prox := proxy_client.ProxyClientResponse{
		ID:           aDTO.ID,
		Title:        aDTO.Title,
		OS:           aDTO.OS,
		DownloadLink: aDTO.DownloadLink,
		CreatedAt:    aDTO.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    aDTO.UpdatedAt.Format(time.RFC3339),
	}

	logger.Instance.Debugf("Mapped app dto to response: %+v", prox)

	return prox
}
