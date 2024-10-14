package client

import (
	"time"

	"github.com/gobuffalo/envy"
)

type Config struct {
	ServerAddr string
	Timeout    time.Duration
}

func NewConfig() *Config {
	cfg := Config{
		ServerAddr: envy.Get("SERVER_ADDR", "localhost:9090"),
		Timeout:    10 * time.Second,
	}
	timeout := envy.Get("TIMEOUT", "10s")
	if duration, err := time.ParseDuration(timeout); err == nil {
		cfg.Timeout = duration
	}

	return &cfg
}
