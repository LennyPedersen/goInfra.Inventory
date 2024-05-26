package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func main() {
    secret := "your-256-bit-secret"

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "client": "client-id",
        "exp":    time.Now().Add(time.Hour * 1).Unix(),
    })

    tokenString, err := token.SignedString([]byte(secret))
    if err != nil {
        fmt.Println("Error generating token:", err)
        return
    }

    fmt.Println("Generated JWT:", tokenString)
}