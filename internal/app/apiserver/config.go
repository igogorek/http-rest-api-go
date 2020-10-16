package apiserver

import "github.com/igogorek/http-rest-api-go/internal/app/store"

// Config type for apiserver
type Config struct {
	BindAddr string `toml:"bind_addr"`
	LogLevel string `toml:"log_level"`
	Store    *store.Config
}

// NewConfig function to initialise config
func NewConfig() *Config {
	return &Config{
		BindAddr: ":8080",
		LogLevel: "debug",
		Store:    store.NewConfig(),
	}
}
