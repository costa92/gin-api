package main

import (
	"context"

	logger2 "github.com/costa92/go-web/pkg/logger"
)

func main() {
	// logger配置
	opts := &logger2.Options{
		Level:            "debug",
		Format:           "console",
		EnableColor:      true,
		DisableCaller:    true,
		OutputPaths:      []string{"test.log", "stdout"},
		ErrorOutputPaths: []string{"error.log"},
	}
	// 初始化全局logger
	logger2.Init(opts)
	defer logger2.Flush()

	lv := logger2.WithValues("X-Request-ID", "7a7b9f24-4cae-4b2a-9464-69088b45b904")

	lv.Infof("Start to call pirntString function")
	ctx := lv.WithContext(context.Background())

	lc := logger2.FromContext(ctx)
	ln := lv.WithName("test")
	ln.Info("123123")
	lc.Infof("Hello %s", "str")
	lc.Infow("对对对", "ddd", "wee")
	lc.Info("Message printed with [WithContext] logger")
}
