package config

import (
	"time"
)

// Server defines the general server configuration.
type Server struct {
	Addr string
	Path string
}

// Target defines the target specific configuration.
type Target struct {
	Address string
	Timeout time.Duration
}

// Config is a combination of all available configurations.
type Config struct {
	Server Server
	Target Target
}

// Load initializes a default configuration struct.
func Load() *Config {
	return &Config{}
}
