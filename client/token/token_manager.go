package token

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Token struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
}

var tokenFilePath = filepath.Join(os.TempDir(), "client_tokens.json")

func SaveTokens(tokens Token) error {
    data, err := json.Marshal(tokens)
    if err != nil {
        return err
    }
    return os.WriteFile(tokenFilePath, data, 0644)
}

func LoadTokens() (Token, error) {
    var tokens Token
    data, err := os.ReadFile(tokenFilePath)
    if err != nil {
        return tokens, err
    }
    if err := json.Unmarshal(data, &tokens); err != nil {
        return tokens, err
    }
    return tokens, nil
}

func SaveRefreshToken(refreshToken string) error {
    return os.WriteFile("refresh_token.txt", []byte(refreshToken), 0644)
}

func LoadRefreshToken() (string, error) {
    refreshToken, err := os.ReadFile("refresh_token.txt")
    if err != nil {
        return "", err
    }
    return string(refreshToken), nil
}