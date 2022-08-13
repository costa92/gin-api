package config

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"github.com/costa92/go-web/internal/option"
)

type Config struct {
	ServerConf  *ServerConf          `mapstructure:"service"`
	MysqlConfig *option.MySQLOptions `mapstructure:"mysql"`
}

type ServerConf struct {
	Name        string   `json:"name" mapstructure:"name"`
	Port        string   `json:"port" mapstructure:"port"`
	Middlewares []string `json:"middlewares" mapstructure:"middlewares"`
}

func NewConfig() (*Config, error) {
	var c Config
	err := viper.Unmarshal(&c)
	if err != nil {
		return nil, err
	}
	log.Info().Msgf("service config:%s", c.ServerConf)
	return &c, nil
}
