package domain

import "time"

type Inventory struct {
    IP               string    `json:"ip"`
    Hostname         string    `json:"hostname"`
    Services         string    `json:"services"`
    OSVersion        string    `json:"os_version"`
    OpenPorts        string    `json:"open_ports"`
    LastReportedDate time.Time `json:"last_reported_date"`
    NextReportDate   time.Time `json:"next_report_date"`
    Health           string    `json:"health"`
}