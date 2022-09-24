package main

import (
	logger2 "github.com/costa92/go-web/pkg/logger"
)

func main() {
	// OptionsTest()
	logger2.Infow("qeqwe", "wqeqwe", "eqweqwe")
}

func OptionsTest() {
	opts := &logger2.Options{
		Level:            "info",
		Format:           "console",
		EnableColor:      true,
		EnableCaller:     true,
		OutputPaths:      []string{"test.log", "stdout"},
		ErrorOutputPaths: []string{},
	}
	logger2.Init(opts)
	defer logger2.Flush()
	logger2.Info("qwee")
}
