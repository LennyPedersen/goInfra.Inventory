package application

import (
	"bytes"
	"encoding/json"
	"fmt"
	"goInfraInventory-client/internal/domain"
	"net/http"
)

type ReportService struct {
    URL string
    JWT string
}

func (s *ReportService) ReportInventory(inventory domain.Inventory) error {
    data, err := json.Marshal(inventory)
    if err != nil {
        return err
    }

    req, err := http.NewRequest("POST", s.URL, bytes.NewBuffer(data))
    if err != nil {
        return err
    }
    req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.JWT))
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode == http.StatusUnauthorized {
        return fmt.Errorf("Unauthorized")
    }

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("server returned non-200 status: %v", resp.Status)
    }

    return nil
}