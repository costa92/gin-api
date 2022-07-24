package config

import (
	"fmt"

	"github.com/costa92/go-web/internal/option"
	"github.com/spf13/viper"
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
	fmt.Println(c.ServerConf)
	return &c, nil
}
