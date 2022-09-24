package options

import (
	"encoding/json"

	cliflag "github.com/costa92/component-base/pkg/cli/flag"
	"github.com/marmotedu/component-base/pkg/util/idutil"

	"github.com/costa92/go-web/internal/option"
	"github.com/costa92/go-web/internal/server"
	"github.com/costa92/go-web/pkg/logger"
)

type Options struct {
	MysqlConfig             *option.MySQLOptions           `json:"mysql"   mapstructure:"mysql"`
	Logger                  *logger.Options                `json:"log"     mapstructure:"log"`
	GenericServerRunOptions *option.ServerRunOptions       `json:"server"  mapstructure:"server"`
	InsecureServingOptions  *option.InsecureServingOptions `json:"insecure"  mapstructure:"insecure"`
	Jwt                     *option.JwtOptions             `json:"jwt"    mapstructure:"jwt"`
	RedisOptions            *option.RedisOptions           `json:"redis"  mapstructure:"redis"`
	SecureServing           *option.SecureServingOptions   `json:"secure" mapstructure:"secure"`
}

func NewOptions() *Options {
	o := Options{
		GenericServerRunOptions: option.NewServerRunOptions(),
		InsecureServingOptions:  option.NewInsecureServingOptions(),
		MysqlConfig:             option.NewMySQLOptions(),
		RedisOptions:            option.NewRedisOptions(),
		Jwt:                     option.NewJwtOptions(),
		Logger:                  logger.NewOptions(),
		SecureServing:           option.NewSecureServingOptions(),
	}
	return &o
}

func (o *Options) ApplyTo(c *server.Config) error {
	return nil
}

func (o *Options) Flags() (fss cliflag.NamedFlagSets) {
	o.GenericServerRunOptions.AddFlags(fss.FlagSet("server"))
	o.Logger.AddFlags(fss.FlagSet("log"))
	o.MysqlConfig.AddFlags(fss.FlagSet("mysql"))
	o.InsecureServingOptions.AddFlags(fss.FlagSet("insecure serving"))
	o.Jwt.AddFlags(fss.FlagSet("jwt"))
	o.RedisOptions.AddFlags(fss.FlagSet("redis"))
	o.SecureServing.AddFlags(fss.FlagSet("secure serving"))
	return fss
}

func (o *Options) String() string {
	data, _ := json.Marshal(o)

	return string(data)
}

func (o *Options) Complete() error {
	if o.Jwt.Key == "" {
		o.Jwt.Key = idutil.NewSecretKey()
	}
	return nil
}
