package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"goInfraInventory-server/internal/application"
	"goInfraInventory-server/internal/infrastructure/persistence"
	httpInterface "goInfraInventory-server/internal/interfaces/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    dbUser := os.Getenv("MYSQL_USER")
    dbPassword := os.Getenv("MYSQL_PASSWORD")
    dbName := os.Getenv("MYSQL_DATABASE")
    dbPort := os.Getenv("MYSQL_PORT")
    serverPort := os.Getenv("SERVER_PORT")

    dsn := fmt.Sprintf("%s:%s@tcp(localhost:%s)/%s", dbUser, dbPassword, dbPort, dbName)

    db, err := sql.Open("mysql", dsn)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    inventoryRepo := &persistence.MySQLInventoryRepository{DB: db}
    inventoryService := &application.InventoryService{Repo: inventoryRepo}
    router := httpInterface.NewRouter(inventoryService)

    log.Printf("Starting server on :%s\n", serverPort)
    log.Fatal(http.ListenAndServe(":"+serverPort, router))
}