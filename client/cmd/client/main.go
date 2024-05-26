package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"goInfraInventory-client/internal/application"
	"goInfraInventory-client/internal/domain"
	"goInfraInventory-client/internal/infrastructure/system"
	"goInfraInventory-client/token"

	"github.com/joho/godotenv"
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	serverURL := os.Getenv("SERVER_URL")
	if serverURL == "" {
		log.Fatalf("SERVER_URL must be provided in the .env file")
	}

	initialAccessToken := os.Getenv("INITIAL_ACCESS_TOKEN")
	if initialAccessToken == "" {
		log.Fatalf("INITIAL_ACCESS_TOKEN must be provided in the .env file")
	}

	tokens, err := token.LoadTokens()
	if err != nil {
		log.Println("No existing tokens found, using initial access token...")
		tokens = token.Token{
			AccessToken: initialAccessToken,
		}

		// Initial authentication to get refresh token
		if err := authenticateAndStoreTokens(serverURL, &tokens); err != nil {
			log.Fatalf("Failed to authenticate and store tokens: %v", err)
		}
	}

	ip, err := system.GetLocalIP()
	if err != nil {
		log.Fatalf("Failed to get local IP: %v", err)
	}

	hostname, err := system.GetHostname()
	if err != nil {
		log.Fatalf("Failed to get hostname: %v", err)
	}

	services, err := system.GetServices()
	if err != nil {
		log.Fatalf("Failed to get services: %v", err)
	}

	osVersion, err := system.GetOSVersion()
	if err != nil {
		log.Fatalf("Failed to get OS version: %v", err)
	}

	openPorts, err := system.GetOpenPorts()
	if err != nil {
		log.Fatalf("Failed to get open ports: %v", err)
	}

	now := time.Now()
	nextReport := now.Add(time.Minute * 15)

	inventory := domain.Inventory{
		IP:               ip,
		Hostname:         hostname,
		Services:         services,
		OSVersion:        osVersion,
		OpenPorts:        openPorts,
		LastReportedDate: now,
		NextReportDate:   nextReport,
		Health:           "Healthy",
	}

	reportService := application.ReportService{
		URL: serverURL + "/inventory",
		JWT: tokens.AccessToken,
	}

	for {
		if err := reportService.ReportInventory(inventory); err != nil {
			if err.Error() == "Unauthorized" {
				tokens, err = renewAccessToken(serverURL)
				if err != nil {
					log.Fatalf("Failed to renew tokens: %v", err)
				}
				reportService.JWT = tokens.AccessToken
				continue
			}
			log.Fatalf("Failed to report inventory: %v", err)
		}
		inventory.LastReportedDate = time.Now()
		inventory.NextReportDate = inventory.LastReportedDate.Add(time.Minute * 15)
		time.Sleep(time.Minute * 15) // Report inventory every 15 minutes
	}
}

func authenticateAndStoreTokens(serverURL string, tokens *token.Token) error {
	client := &http.Client{}
	req, err := http.NewRequest("POST", serverURL+"/authenticate", nil) // Assuming an endpoint that provides the initial tokens
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to authenticate: %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	tokens.AccessToken = tokenResp.AccessToken
	tokens.RefreshToken = tokenResp.RefreshToken

	if err := token.SaveTokens(*tokens); err != nil {
		return fmt.Errorf("failed to save tokens: %w", err)
	}

	if err := token.SaveRefreshToken(tokens.RefreshToken); err != nil {
		return fmt.Errorf("failed to save refresh token: %w", err)
	}

	return nil
}

func renewAccessToken(serverURL string) (token.Token, error) {
	refreshToken, err := token.LoadRefreshToken()
	if err != nil {
		return token.Token{}, fmt.Errorf("failed to read refresh token: %w", err)
	}

	client := &http.Client{}
	reqBody := fmt.Sprintf(`{"refresh_token":"%s"}`, refreshToken)
	req, reqErr := http.NewRequest("POST", serverURL+"/refresh", bytes.NewBuffer([]byte(reqBody)))
	if reqErr != nil {
		return token.Token{}, fmt.Errorf("failed to create request: %w", reqErr)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return token.Token{}, fmt.Errorf("failed to refresh token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return token.Token{}, fmt.Errorf("failed to refresh token: %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return token.Token{}, fmt.Errorf("failed to read response: %w", err)
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return token.Token{}, fmt.Errorf("failed to parse response: %w", err)
	}

	tokens := token.Token{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
	}

	if err := token.SaveTokens(tokens); err != nil {
		return token.Token{}, fmt.Errorf("failed to save tokens: %w", err)
	}

	return tokens, nil
}