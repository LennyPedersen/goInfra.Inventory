package http

import (
	"goInfraInventory-server/internal/application"
	"goInfraInventory-server/internal/infrastructure/http"

	"github.com/gorilla/mux"
)

func NewRouter(service *application.InventoryService) *mux.Router {
    handler := &http.InventoryHandler{
        Service: service,
    }

    router := mux.NewRouter()
    router.HandleFunc("/inventory", handler.ReportInventory).Methods("POST")
    router.Use(JWTMiddleware)

    return router
}