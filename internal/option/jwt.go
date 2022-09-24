package option

import (
	"fmt"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/spf13/pflag"

	"github.com/costa92/go-web/internal/server"
)

type JwtOptions struct {
	Realm      string        `json:"realm"       mapstructure:"realm"`
	Key        string        `json:"key"         mapstructure:"key"`
	Timeout    time.Duration `json:"timeout"     mapstructure:"timeout"`
	MaxRefresh time.Duration `json:"max-refresh" mapstructure:"max-refresh"`
}

// NewJwtOptions 实例化获取默认值
func NewJwtOptions() *JwtOptions {
	defaults := server.NewConfig()
	return &JwtOptions{
		Realm:      defaults.Jwt.Realm,
		Key:        defaults.Jwt.Key,
		Timeout:    defaults.Jwt.Timeout,
		MaxRefresh: defaults.Jwt.MaxRefresh,
	}
}

// ApplyTo 获取的配置
func (s *JwtOptions) ApplyTo(c *server.Config) error {
	c.Jwt = &server.JwtInfo{
		Realm:      s.Realm,
		Key:        s.Key,
		Timeout:    s.Timeout,
		MaxRefresh: s.MaxRefresh,
	}
	return nil
}

func (s *JwtOptions) Validate() []error {
	var errs []error
	if !govalidator.StringLength(s.Key, "6", "32") {
		errs = append(errs, fmt.Errorf("--secret-key must larger than 5 and little than 33"))
	}
	return errs
}

func (s *JwtOptions) AddFlags(fs *pflag.FlagSet) {
	if fs == nil {
		return
	}
	fs.StringVar(&s.Realm, "jwt.realm", s.Realm, "Realm name to display to the user.")
	fs.StringVar(&s.Key, "jwt.key", s.Key, "Private key used to sign jwt token.")
	fs.DurationVar(&s.Timeout, "jwt.timeout", s.Timeout, "JWT token timeout.")
	fs.DurationVar(&s.MaxRefresh, "jwt.max-refresh", s.MaxRefresh, ""+
		"This field allows clients to refresh their token until MaxRefresh has passed.")
}
