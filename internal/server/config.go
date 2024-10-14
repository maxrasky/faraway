package server

import (
	"fmt"

	"github.com/gobuffalo/envy"
)

type Config struct {
	Port string
}

func NewConfig() *Config {
	port := envy.Get("PORT", "9090")

	return &Config{
		Port: fmt.Sprintf(":%s", port),
	}
}
