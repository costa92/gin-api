package option

import "time"

type MySQLOptions struct {
	Addr                  string        `json:"addr" mapstructure:"addr"`
	User                  string        `json:"user" mapstructure:"user"`
	Pass                  string        `json:"pass" mapstructure:"pass"`
	Database              string        `json:"database" mapstructure:"database"`
	MaxIdleConnections    int           `json:"max-idle-connections" mapstructure:"max_idle_connections"`
	MaxOpenConnections    int           `json:"max-open-connections" mapstructure:"max-open-connections"`
	MaxConnectionLeftTime time.Duration `json:"max-connection-left-time" mapstructure:"max-connection-left-time"`
	LogLevel              int           `json:"log-level" mapstructure:"log-level"`
}
