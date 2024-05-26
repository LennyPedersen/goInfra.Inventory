package http

import (
	"encoding/json"
	"goInfraInventory-server/internal/application"
	"goInfraInventory-server/internal/domain"
	"net/http"
	"time"
)

type InventoryHandler struct {
    Service *application.InventoryService
}

func (h *InventoryHandler) ReportInventory(w http.ResponseWriter, r *http.Request) {
    var inventory domain.Inventory
    if err := json.NewDecoder(r.Body).Decode(&inventory); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Update last reported date and next report date
    inventory.LastReportedDate = time.Now()
    inventory.NextReportDate = inventory.LastReportedDate.Add(time.Minute * 15)

    if err := h.Service.ReportInventory(inventory); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}