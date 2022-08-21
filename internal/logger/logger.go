package logger

import (
	"context"
	"fmt"
	"log"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/costa92/go-web/internal/logger/klog"
)

type InfoLogger interface {
	Info(msg string, fields ...Field)
	Infof(format string, v ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	Enabled() bool
}

type Logger interface {
	InfoLogger
	Error(msg string, fields ...Field)
	Errorf(format string, v ...interface{})
	V(level Level) InfoLogger
	Write(p []byte) (n int, err error)
	WithValues(keysAndValues ...interface{}) Logger
	WithName(name string) Logger
	WithContext(ctx context.Context) context.Context
	Flush()
}

var _ Logger = &zapLogger{}

type noopInfoLogger struct{}

func (l *noopInfoLogger) Enabled() bool                    { return false }
func (l *noopInfoLogger) Info(_ string, _ ...Field)        {}
func (l *noopInfoLogger) Infof(_ string, _ ...interface{}) {}
func (l *noopInfoLogger) Infow(_ string, _ ...interface{}) {}

var disabledInfoLogger = &noopInfoLogger{}

type infoLogger struct {
	level Level
	log   *zap.Logger
}

func (l *infoLogger) Enabled() bool { return true }
func (l *infoLogger) Info(msg string, fields ...Field) {
	if checkedEntry := l.log.Check(l.level, msg); checkedEntry != nil {
		checkedEntry.Write(fields...)
	}
}

func (l *infoLogger) Infow(msg string, keysAndValues ...interface{}) {
	if checkedEntry := l.log.Check(l.level, msg); checkedEntry != nil {
		checkedEntry.Write(handleFields(l.log, keysAndValues)...)
	}
}

func (l *infoLogger) Infof(format string, args ...interface{}) {
	if checkedEntry := l.log.Check(l.level, fmt.Sprintf(format, args...)); checkedEntry != nil {
		checkedEntry.Write()
	}
}

// zapLogger is a logr.Logger that uses Zap to log.
type zapLogger struct {
	// NB: this looks very similar to zap.SugaredLogger, but
	// deals with our desire to have multiple verbosity levels.
	zapLogger *zap.Logger
	infoLogger
}

func handleFields(l *zap.Logger, args []interface{}, additional ...zap.Field) []zap.Field {
	// a slightly modified version of zap.SugaredLogger.sweetenFields
	if len(args) == 0 {
		// fast-return if we have no suggared fields.
		return additional
	}
	// unlike Zap, we can be pretty sure users aren't passing structured
	// fields (since logr has no concept of that), so guess that we need a
	// little less space.
	fields := make([]zap.Field, 0, len(args)/2+len(additional))
	for i := 0; i < len(args); {
		// check just in case for strongly-typed Zap fields, which is illegal (since
		// it breaks implementation agnosticism), so we can give a better error message.
		if _, ok := args[i].(zap.Field); ok {
			l.DPanic("strongly-typed Zap Field passed to logr", zap.Any("zap field", args[i]))
			break
		}
		// make sure this isn't a mismatched key
		if i == len(args)-1 {
			l.DPanic("odd number of arguments passed as key-value pairs for logging", zap.Any("ignored key", args[i]))
			break
		}
		// process a key-value pair,
		// ensuring that the key is a string
		key, val := args[i], args[i+1]
		keyStr, isString := key.(string)
		if !isString {
			// if the key isn't a string, DPanic and stop logging
			l.DPanic(
				"non-string key argument passed to logging, ignoring all later arguments",
				zap.Any("invalid key", key),
			)
			break
		}
		fields = append(fields, zap.Any(keyStr, val))
		i += 2
	}
	return append(fields, additional...)
}

var logger *zapLogger

var (
	std = New(NewOptions())
	mu  sync.Mutex
)

func Init(opts *Options) {
	mu.Lock()
	defer mu.Unlock()
	std = New(opts)
}

func New(opts *Options) *zapLogger {
	if opts == nil {
		opts = NewOptions()
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     timeEncoder,
		EncodeDuration: milliSecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	// when output to local path, with color is forbidden
	if opts.EnableColor {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(opts.Level)); err != nil {
		zapLevel = InfoLevel
	}

	loggerConfig := &zap.Config{
		Level:             zap.NewAtomicLevelAt(zapLevel),
		Development:       false,
		DisableCaller:     !opts.EnableCaller,
		DisableStacktrace: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         opts.Format,
		EncoderConfig:    encoderConfig,
		OutputPaths:      opts.OutputPaths,
		ErrorOutputPaths: opts.ErrorOutputPaths,
	}

	var err error
	l, err := loggerConfig.Build(zap.AddStacktrace(zapcore.PanicLevel), zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}
	logger = &zapLogger{
		zapLogger: l,
		infoLogger: infoLogger{
			log:   l,
			level: zap.InfoLevel,
		},
	}
	klog.InitLogger(l)
	zap.RedirectStdLog(l)

	return logger
}

func SugaredLogger() *zap.SugaredLogger {
	return std.zapLogger.Sugar()
}

// StdInfoLogger returns logger of standard library which writes to supplied zap
// logger at info level.
func StdInfoLogger() *log.Logger {
	if logger == nil {
		return nil
	}
	if l, err := zap.NewStdLogAt(logger.zapLogger, zapcore.InfoLevel); err == nil {
		return l
	}
	return nil
}

func StdErrLogger() *log.Logger {
	if std == nil {
		return nil
	}
	if l, err := zap.NewStdLogAt(logger.zapLogger, zapcore.ErrorLevel); err == nil {
		return l
	}
	return nil
}
func V(level Level) InfoLogger { return std.V(level) }

func (l *zapLogger) V(level Level) InfoLogger {
	if l.zapLogger.Core().Enabled(level) {
		return &infoLogger{
			level: level,
			log:   l.zapLogger,
		}
	}
	return disabledInfoLogger
}

func (l *zapLogger) Write(p []byte) (n int, err error) {
	l.zapLogger.Info(string(p))
	return len(p), nil
}

func WithValues(keysAndValues ...interface{}) Logger {
	return logger.WithValues(keysAndValues...)
}

func (l *zapLogger) WithValues(keysAndValues ...interface{}) Logger {
	newLogger := l.zapLogger.With(handleFields(l.zapLogger, keysAndValues)...)
	return NewLogger(newLogger)
}

func WithName(s string) Logger { return logger.WithName(s) }

func (l *zapLogger) WithName(name string) Logger {
	newLogger := l.zapLogger.Named(name)
	return NewLogger(newLogger)
}

func Flush() { std.Flush() }

func (l *zapLogger) Flush() {
	_ = l.zapLogger.Sync()
}

func NewLogger(l *zap.Logger) Logger {
	return &zapLogger{
		zapLogger: l,
		infoLogger: infoLogger{
			log:   l,
			level: zap.InfoLevel,
		},
	}
}

func ZapLogger() *zap.Logger {
	return logger.zapLogger
}

func CheckIntLevel(level int32) bool {
	var lvl zapcore.Level
	if level < 5 {
		lvl = zapcore.InfoLevel
	} else {
		lvl = zapcore.DebugLevel
	}
	checkEntry := std.zapLogger.Check(lvl, "")

	return checkEntry != nil
}
