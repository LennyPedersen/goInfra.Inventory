package persistence

import (
	"database/sql"
	"fmt"
	"goInfraInventory-server/internal/domain"
	"os"
	"time"
)

type MySQLInventoryRepository struct {
    DB *sql.DB
}

func (r *MySQLInventoryRepository) Save(inventory domain.Inventory) error {
    tableName := os.Getenv("INVENTORY_TABLE")
    if tableName == "" {
        return fmt.Errorf("table name not set in environment variables")
    }

    query := fmt.Sprintf("INSERT INTO %s (ip, hostname, services, os_version, open_ports, last_reported_date, next_report_date, health) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", tableName)
    _, err := r.DB.Exec(query, inventory.IP, inventory.Hostname, inventory.Services, inventory.OSVersion, inventory.OpenPorts, inventory.LastReportedDate, inventory.NextReportDate, inventory.Health)
    return err
}

func (r *MySQLInventoryRepository) UpdateHealth(ip string, health string, nextReportDate time.Time) error {
    tableName := os.Getenv("INVENTORY_TABLE")
    if tableName == "" {
        return fmt.Errorf("table name not set in environment variables")
    }

    query := fmt.Sprintf("UPDATE %s SET health = ?, next_report_date = ? WHERE ip = ?", tableName)
    _, err := r.DB.Exec(query, health, nextReportDate, ip)
    return err
}