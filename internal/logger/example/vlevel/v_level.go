package main

import "github.com/costa92/go-web/internal/logger"

func main() {
	defer logger.Flush()
	logger.V(1).Info("his is a V level message")
}
