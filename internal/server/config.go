package server

import (
	"fmt"
	"time"

	"github.com/gobuffalo/envy"
)

type Config struct {
	Port     string
	Deadline time.Duration
}

func NewConfig() *Config {
	port := envy.Get("PORT", "9090")
	deadline := envy.Get("DEADLINE", "10s")

	cfg := Config{
		Port:     fmt.Sprintf(":%s", port),
		Deadline: 10 * time.Second,
	}
	if deadlineTimeout, err := time.ParseDuration(deadline); err == nil {
		cfg.Deadline = deadlineTimeout
	}

	return &cfg
}
