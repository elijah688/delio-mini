package config

import (
	"os"

	"github.com/Finnhub-Stock-API/finnhub-go/v2"
)

const (
	FH_TOKEN = "FH_TOKEN"
)

// Generate Config for the FH Service
func New() *finnhub.Configuration {
	cfg, auth := finnhub.NewConfiguration(), os.Getenv(FH_TOKEN)
	if auth == "" {
		panic("finnhub authentication token not set")
	}
	cfg.AddDefaultHeader("X-Finnhub-Token", auth)

	return cfg

}
