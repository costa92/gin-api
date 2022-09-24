package server

import (
	"net"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Config struct {
	Mode            string
	Middlewares     []string
	SecureServing   *SecureServingInfo
	InsecureServing *InsecureServingInfo
	Healthz         bool
	Jwt             *JwtInfo
}

// JwtInfo defines jwt fields used to create jwt authentication middleware.
type JwtInfo struct {
	// defaults to "iam jwt"
	Realm string
	// defaults to empty
	Key string
	// defaults to one hour
	Timeout time.Duration
	// defaults to zero
	MaxRefresh time.Duration
}

// CertKey contains configuration items related to certificate.
type CertKey struct {
	// CertFile is a file containing a PEM-encoded certificate, and possibly the complete certificate chain
	CertFile string
	// KeyFile is a file containing a PEM-encoded private key for the certificate specified by CertFile
	KeyFile string
}

type SecureServingInfo struct {
	BindAddress string
	BindPort    int
	CertKey     CertKey
}

// Address join host IP address and host port number into a address string, like: 0.0.0.0:8443.
func (s *SecureServingInfo) Address() string {
	return net.JoinHostPort(s.BindAddress, strconv.Itoa(s.BindPort))
}

type InsecureServingInfo struct {
	Address string
}

func NewConfig() *Config {
	return &Config{
		Healthz:     true,
		Mode:        gin.ReleaseMode,
		Middlewares: []string{},
		Jwt: &JwtInfo{
			Realm:      "iam jwt",
			Timeout:    1 * time.Hour,
			MaxRefresh: 1 * time.Hour,
		},
	}
}

type CompletedConfig struct {
	*Config
}

func (c *Config) Complete() CompletedConfig {
	return CompletedConfig{c}
}

func (c CompletedConfig) New() (*GenericAPIServer, error) {
	// setMode before gin.New()
	gin.SetMode(c.Mode)
	s := &GenericAPIServer{
		healthz:             c.Healthz,
		middlewares:         c.Middlewares,
		Engine:              gin.New(),
		SecureServingInfo:   c.SecureServing,
		InsecureServingInfo: c.InsecureServing,
	}
	initGenericAPIServer(s)
	return s, nil
}
