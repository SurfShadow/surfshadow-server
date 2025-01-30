package routes

import (
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/SurfShadow/surfshadow-server/internal/presentation/handlers"
	"github.com/SurfShadow/surfshadow-server/internal/presentation/middleware"
)

func InitRoutes(proxyClientHandler *handlers.ProxyClientHandler) *mux.Router {
	router := mux.NewRouter()

	router.Use(middleware.LoggingMiddleware)

	apiRouter := router.PathPrefix("/api/v1").Subrouter()

	apiRouter.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	apiRouter.HandleFunc("/proxy-clients", proxyClientHandler.GetAllProxyClients).Methods("GET")
	apiRouter.HandleFunc("/proxy-clients", proxyClientHandler.CreateProxyClient).Methods("POST")
	apiRouter.HandleFunc("/clients/{vpn_client_id}", proxyClientHandler.UpdateProxyClient).Methods("PATCH")
	apiRouter.HandleFunc("/clients/{vpn_client_id}", proxyClientHandler.DeleteProxyClient).Methods("DELETE")

	return router
}
