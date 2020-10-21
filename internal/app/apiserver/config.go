package apiserver

// Config type for apiserver
type Config struct {
	BindAddr    string `toml:"bind_addr"`
	LogLevel    string `toml:"log_level"`
	DatabaseURL string `toml:"database_url"`
}

// NewConfig function to initialise config
func NewConfig() *Config {
	return &Config{
		BindAddr:    ":8080",
		LogLevel:    "debug",
		DatabaseURL: "",
	}
}
