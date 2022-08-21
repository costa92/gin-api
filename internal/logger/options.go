package logger

import (
	"encoding/json"
	"fmt"
	"strings"

	"go.uber.org/zap/zapcore"
)

const (
	flagLevel            = "log.level"
	flagFormat           = "log.format"
	flagEnableColor      = "log.enable-color"
	flagEnableCaller     = "log.enable-caller"
	flagOutputPaths      = "log.output-paths"
	flagErrorOutputPaths = "log.error-output-paths"

	consoleFormat = "console"
	jsonFormat    = "json"
)

type Options struct {
	Level             string   `json:"level" yaml:"level" mapstructure:"level"`
	Format            string   `json:"format" yaml:"format" mapstructure:"format"`
	EnableColor       bool     `json:"enable-color" yaml:"enable-color" mapstructure:"enable-color"`
	EnableCaller      bool     `json:"enable-caller" yaml:"enable-caller" mapstructure:"enable-caller"`
	OutputPaths       []string `json:"output-paths" yaml:"output-paths" mapstructure:"output-paths"`
	ErrorOutputPaths  []string `json:"error-output-paths" yaml:"error-output-paths" mapstructure:"error-output-paths"`
	Development       bool     `json:"development"    yaml:"development"    mapstructure:"development"`
	Name              string   `json:"name" yaml:"name"  mapstructure:"name"`
	DisableCaller     bool     `json:"disable-caller"  yaml:"disable-caller"   mapstructure:"disable-caller"`
	DisableStacktrace bool     `json:"disable-stacktrace" yaml:"disable-stacktrace" mapstructure:"disable-stacktrace"`
}

func NewOptions() *Options {
	return &Options{
		Level:            zapcore.InfoLevel.String(),
		Format:           jsonFormat,
		EnableColor:      false,
		EnableCaller:     false,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
}

func (o *Options) Validate() []error {
	var errs []error
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(o.Level)); err != nil {
		errs = append(errs, err)
	}

	format := strings.ToLower(o.Format)
	if format != consoleFormat && format != jsonFormat {
		errs = append(errs, fmt.Errorf("not a valid log format: %q", o.Format))
	}
	return errs
}

func (o *Options) String() string {
	data, _ := json.Marshal(o)
	return string(data)
}
