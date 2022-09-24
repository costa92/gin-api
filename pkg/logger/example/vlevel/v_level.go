package main

import (
	"github.com/costa92/go-web/pkg/logger"
)

func main() {
	defer logger.Flush()
	logger.V(1).Info("his is a V level message")
}
