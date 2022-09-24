package logger

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/pflag"
	"go.uber.org/zap/zapcore"
)

const (
	flagLevel             = "log.level"
	flagFormat            = "log.format"
	flagEnableColor       = "log.enable-color"
	flagEnableCaller      = "log.enable-caller"
	flagDisableStacktrace = "log.disable-stacktrace"
	flagOutputPaths       = "log.output-paths"
	flagErrorOutputPaths  = "log.error-output-paths"
	flagDisableCaller     = "log.disable-caller"
	consoleFormat         = "console"
	jsonFormat            = "json"
	flagDevelopment       = "log.development"
	flagName              = "log.name"
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

func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Level, flagLevel, o.Level, "Minimum log output `LEVEL`.")
	fs.BoolVar(&o.DisableCaller, flagDisableCaller, o.DisableCaller, "Disable output of caller information in the log.")
	fs.BoolVar(&o.DisableStacktrace, flagDisableStacktrace,
		o.DisableStacktrace, "Disable the log to record a stack trace for all messages at or above panic level.")
	fs.StringVar(&o.Format, flagFormat, o.Format, "Log output `FORMAT`, support plain or json format.")
	fs.BoolVar(&o.EnableColor, flagEnableColor, o.EnableColor, "Enable output ansi colors in plain format logs.")
	fs.StringSliceVar(&o.OutputPaths, flagOutputPaths, o.OutputPaths, "Output paths of log.")
	fs.StringSliceVar(&o.ErrorOutputPaths, flagErrorOutputPaths, o.ErrorOutputPaths, "Error output paths of log.")
	fs.BoolVar(
		&o.Development,
		flagDevelopment,
		o.Development,
		"Development puts the logger in development mode, which changes "+
			"the behavior of DPanicLevel and takes stacktraces more liberally.",
	)
	fs.StringVar(&o.Name, flagName, o.Name, "The name of the logger.")
}
