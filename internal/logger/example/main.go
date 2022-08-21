package main

import "github.com/costa92/go-web/internal/logger"

func main() {
	// OptionsTest()
	logger.Infow("qeqwe", "wqeqwe", "eqweqwe")
}

func OptionsTest() {
	opts := &logger.Options{
		Level:            "info",
		Format:           "console",
		EnableColor:      true,
		EnableCaller:     true,
		OutputPaths:      []string{"test.log", "stdout"},
		ErrorOutputPaths: []string{},
	}
	logger.Init(opts)
	defer logger.Flush()
	logger.Info("qwee")
}
