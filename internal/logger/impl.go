package logger

func Debug(msg string, fields ...Field) {
	logger.zapLogger.Debug(msg, fields...)
}

// Debugf method output debug level log.
func Debugf(format string, v ...interface{}) {
	logger.zapLogger.Sugar().Debugf(format, v...)
}

// Info method output info level log.
func Info(msg string, fields ...Field) {
	logger.zapLogger.Info(msg, fields...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	logger.zapLogger.Sugar().Infow(msg, keysAndValues...)
}

// Error method output error level log.
func Error(msg string, fields ...Field) {
	logger.zapLogger.Error(msg, fields...)
}

// Errorf method output error level log.
func Errorf(format string, v ...interface{}) {
	logger.zapLogger.Sugar().Errorf(format, v...)
}
