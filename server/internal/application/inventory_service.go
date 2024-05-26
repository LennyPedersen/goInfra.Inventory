package application

import (
	"goInfraInventory-server/internal/domain"
	"goInfraInventory-server/internal/infrastructure/persistence"
	"time"
)

type InventoryService struct {
    Repo *persistence.MySQLInventoryRepository
}

func (s *InventoryService) ReportInventory(inventory domain.Inventory) error {
    if err := s.Repo.Save(inventory); err != nil {
        return err
    }

    health := "Healthy"
    if time.Since(inventory.LastReportedDate) > 2*time.Minute {
        health = "Unhealthy"
    }

    nextReportDate := time.Now().Add(time.Minute * 15)
    return s.Repo.UpdateHealth(inventory.IP, health, nextReportDate)
}