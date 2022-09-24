package config

import (
	"github.com/costa92/go-web/internal/options"
)

type Config struct {
	*options.Options
}

func CreateConfigFromOptions(opts *options.Options) (*Config, error) {
	return &Config{opts}, nil
}
