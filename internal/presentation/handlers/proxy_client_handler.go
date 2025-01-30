package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"

	"github.com/SurfShadow/surfshadow-server/internal/application/usecases"
	"github.com/SurfShadow/surfshadow-server/internal/presentation/dto/proxy_client"
	"github.com/SurfShadow/surfshadow-server/internal/presentation/mapper"
	"github.com/SurfShadow/surfshadow-server/pkg/logger"
)

const (
	errInvalidRequestBody        = "invalid request body"
	errValidationFailed          = "validation failed"
	errFailedToCreateProxyClient = "failed to create proxy client"
	errFailedToFetchProxyClients = "failed to fetch proxy clients"
	errNoVPNClientsFound         = "no vpn clients found"
	errInvalidVPNClientID        = "invalid vpn client id"
	errFailedToUpdateProxyClient = "failed to update proxy client"
	errFailedToDeleteProxyClient = "failed to delete proxy client"
	errFailedToEncodeResponse    = "failed to encode response"
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

// CreateProxyClient godoc
// @Summary Add a new VPN client
// @Description Add a new VPN client to the database.
// @Tags clients
// @Accept json
// @Produce json
// @Param client body proxy_client.ProxyClientRequest true "New Client"
// @Success 201 {object} proxy_client.ProxyClientResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /proxy-clients [post].
func (h *ProxyClientHandler) CreateProxyClient(w http.ResponseWriter, r *http.Request) {
	var req proxy_client.ProxyClientRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Instance.Errorf("%s: %v", errInvalidRequestBody, err)
		http.Error(w, errInvalidRequestBody, http.StatusBadRequest)

		return
	}

	if err := h.validate.Struct(req); err != nil {
		logger.Instance.Errorf("%s: %v", errValidationFailed, err)
		http.Error(w, errValidationFailed+": "+err.Error(), http.StatusBadRequest)

		return
	}

	appDTO := mapper.MapRequestToAppDTO(req)

	createdClient, err := h.useCase.CreateProxyClient(&appDTO)
	if err != nil {
		logger.Instance.Errorf("%s: %v", errFailedToCreateProxyClient, err)
		http.Error(w, errFailedToCreateProxyClient, http.StatusInternalServerError)

		return
	}

	response := mapper.MapAppDTOToResponse(*createdClient)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		logger.Instance.Errorf("%s: %v", errFailedToEncodeResponse, err)
		http.Error(w, errFailedToEncodeResponse, http.StatusInternalServerError)

		return
	}
}

// GetAllProxyClients godoc
// @Summary Get a list of VPN clients
// @Description Retrieve a list of VPN clients available for download.
// @Tags clients
// @Accept json
// @Produce json
// @Param id query int false "Client ID"
// @Param title query string false "Client Title"
// @Param os query string false "Client OS"
// @Success 200 {array} proxy_client.ProxyClientResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /proxy-clients [get].
func (h *ProxyClientHandler) GetAllProxyClients(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	idParam := queryParams.Get("id")
	titleParam := queryParams.Get("title")
	osParam := queryParams.Get("os")

	clients, err := h.useCase.GetAllProxyClients()
	if err != nil {
		logger.Instance.Errorf("%s: %v", errFailedToFetchProxyClients, err)
		http.Error(w, errFailedToFetchProxyClients, http.StatusInternalServerError)

		return
	}

	var filteredClients []proxy_client.ProxyClientResponse
	for _, client := range clients {
		if (idParam == "" || strconv.FormatInt(client.ID, 10) == idParam) &&
			(titleParam == "" || client.Title == titleParam) &&
			(osParam == "" || client.OS == osParam) {
			filteredClients = append(filteredClients, mapper.MapAppDTOToResponse(*client))
		}
	}

	if len(filteredClients) == 0 {
		http.Error(w, errNoVPNClientsFound, http.StatusNotFound)
		return
	}

	for _, client := range filteredClients {
		logger.Instance.Infof("Filtered client: %+v", client)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(filteredClients)
	if err != nil {
		logger.Instance.Errorf("%s: %v", errFailedToEncodeResponse, err)
		http.Error(w, errFailedToEncodeResponse, http.StatusInternalServerError)

		return
	}
}

// UpdateProxyClient godoc
// @Summary Update an existing VPN client by id
// @Description Update details of an existing VPN client.
// @Tags clients
// @Accept json
// @Produce json
// @Param vpn_client_id path int true "Client ID"
// @Param client body proxy_client.ProxyClientRequest true "Update Client"
// @Success 200 {object} map[string]bool
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /clients/{vpn_client_id} [patch].
func (h *ProxyClientHandler) UpdateProxyClient(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.ParseInt(vars["vpn_client_id"], 10, 64)
	if err != nil {
		logger.Instance.Errorf("%s: %v", errInvalidVPNClientID, err)
		http.Error(w, errInvalidVPNClientID, http.StatusBadRequest)

		return
	}

	var req proxy_client.ProxyClientRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Instance.Errorf("%s: %v", errInvalidRequestBody, err)
		http.Error(w, errInvalidRequestBody, http.StatusBadRequest)

		return
	}

	if err = h.validate.Struct(req); err != nil {
		logger.Instance.Errorf("%s: %v", errValidationFailed, err)
		http.Error(w, errValidationFailed+": "+err.Error(), http.StatusBadRequest)

		return
	}

	appDTO := mapper.MapRequestToAppDTO(req)
	appDTO.ID = id

	err = h.useCase.UpdateProxyClient(&appDTO)
	if err != nil {
		logger.Instance.Errorf("%s: %v", errFailedToUpdateProxyClient, err)
		http.Error(w, errFailedToUpdateProxyClient, http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(map[string]bool{"success": true})
	if err != nil {
		logger.Instance.Errorf("%s: %v", errFailedToEncodeResponse, err)
		http.Error(w, errFailedToEncodeResponse, http.StatusInternalServerError)

		return
	}
}

// DeleteProxyClient godoc
// @Summary Delete a VPN client by id
// @Description Remove a VPN client from the database by its unique id.
// @Tags clients
// @Param vpn_client_id path int true "Client ID"
// @Success 204
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /clients/{vpn_client_id} [delete].
func (h *ProxyClientHandler) DeleteProxyClient(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.ParseInt(vars["vpn_client_id"], 10, 64)
	if err != nil {
		logger.Instance.Errorf("%s: %v", errInvalidVPNClientID, err)
		http.Error(w, errInvalidVPNClientID, http.StatusBadRequest)

		return
	}

	err = h.useCase.DeleteProxyClient(id)
	if err != nil {
		logger.Instance.Errorf("%s: %v", errFailedToDeleteProxyClient, err)
		http.Error(w, errFailedToDeleteProxyClient, http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
