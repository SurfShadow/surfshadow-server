package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/SurfShadow/surfshadow-server/internal/application/usecases"
	"github.com/SurfShadow/surfshadow-server/internal/presentation/dto"
	"github.com/SurfShadow/surfshadow-server/internal/presentation/mapper"
)

type ProxyClientHandler struct {
	useCase  *usecases.ProxyClientUseCase
	validate *validator.Validate
}

func NewProxyClientHandler(useCase *usecases.ProxyClientUseCase) *ProxyClientHandler {
	return &ProxyClientHandler{
		useCase:  useCase,
		validate: validator.New(),
	}
}

func (h *ProxyClientHandler) CreateProxyClient(w http.ResponseWriter, r *http.Request) {
	var req dto.ProxyClientRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	appDTO := mapper.MapRequestToAppDTO(req)

	createdClient, err := h.useCase.CreateProxyClient(&appDTO)
	if err != nil {
		http.Error(w, "Failed to create proxy client", http.StatusInternalServerError)
		return
	}

	response := mapper.MapAppDTOToResponse(*createdClient)

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *ProxyClientHandler) GetAllProxyClients(w http.ResponseWriter, r *http.Request) {
	clients, err := h.useCase.GetAllProxyClients()
	if err != nil {
		http.Error(w, "Failed to fetch proxy clients", http.StatusInternalServerError)
		return
	}

	var responses = make([]dto.ProxyClientResponse, 0)
	for _, client := range clients {
		responses = append(responses, mapper.MapAppDTOToResponse(*client))
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(responses)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
