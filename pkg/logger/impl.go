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

func Infof(format string, v ...interface{}) {
	logger.zapLogger.Sugar().Infof(format, v...)
}

// Error method output error level log.
func Error(msg string, fields ...Field) {
	logger.zapLogger.Error(msg, fields...)
}

// Errorf method output error level log.
func Errorf(format string, v ...interface{}) {
	logger.zapLogger.Sugar().Errorf(format, v...)
}

// Errorw
func Errorw(msg string, keysAndValues ...interface{}) {
	logger.zapLogger.Sugar().Errorw(msg, keysAndValues...)
}

// Warn method output warning level log.
func Warn(msg string, fields ...Field) {
	logger.zapLogger.Warn(msg, fields...)
}

// Warnf method output warning level log.
func Warnf(format string, v ...interface{}) {
	logger.zapLogger.Sugar().Warnf(format, v...)
}

// Warnw method output warning level log.
func Warnw(msg string, keysAndValues ...interface{}) {
	logger.zapLogger.Sugar().Warnw(msg, keysAndValues...)
}

// Panic method output panic level log and shutdown application.
func Panic(msg string, fields ...Field) {
	logger.zapLogger.Panic(msg, fields...)
}

// Panicf method output panic level log and shutdown application.
func Panicf(format string, v ...interface{}) {
	logger.zapLogger.Sugar().Panicf(format, v...)
}

// Panicw method output panic level log.
func Panicw(msg string, keysAndValues ...interface{}) {
	logger.zapLogger.Sugar().Panicw(msg, keysAndValues...)
}

// Fatalf method output fatal level log.
func Fatalf(format string, v ...interface{}) {
	logger.zapLogger.Sugar().Fatalf(format, v...)
}

// Fatalw method output Fatalw level log.
func Fatalw(msg string, keysAndValues ...interface{}) {
	logger.zapLogger.Sugar().Fatalw(msg, keysAndValues...)
}
