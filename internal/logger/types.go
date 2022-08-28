package logger

import "go.uber.org/zap/zapcore"

const (
	KeyRequestID string = "requestID"
)

type Field = zapcore.Field

// Level is an alias for the level structure in the underlying log frame.
type Level = zapcore.Level

var InfoLevel = zapcore.InfoLevel
