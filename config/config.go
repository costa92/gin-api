package config

import (
	"github.com/spf13/viper"

	"github.com/costa92/go-web/internal/logger"
	"github.com/costa92/go-web/internal/option"
)

type Config struct {
	ServerConf  *ServerConf          `mapstructure:"service"`
	MysqlConfig *option.MySQLOptions `mapstructure:"mysql"`
	Logger      *logger.Options      `mapstructure:"log"`
}

type ServerConf struct {
	Name        string   `json:"name" mapstructure:"name"`
	Port        string   `json:"port" mapstructure:"port"`
	Mode        string   `json:"mode" mapstructure:"mode"`
	Middlewares []string `json:"middlewares" mapstructure:"middlewares"`
}

func NewConfig() (*Config, error) {
	var c Config
	err := viper.Unmarshal(&c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
