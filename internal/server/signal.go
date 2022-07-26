package server

import (
	"os"
	"os/signal"
	"syscall"
)

var onlyOneSignalHandler = make(chan struct{})

var shutdownHandler chan os.Signal

var shutdownSignals = []os.Signal{os.Interrupt, syscall.SIGTERM}

func SetupSignalHandler() <-chan struct{} {
	close(onlyOneSignalHandler)

	shutdownHandler = make(chan os.Signal, 2)

	stop := make(chan struct{})

	signal.Notify(shutdownHandler, shutdownSignals...)

	go func() {
		<-shutdownHandler
		close(stop)
		<-shutdownHandler
		os.Exit(1) // second signal. Exit directly.
	}()
	return stop
}

func RequestShutdown() bool {
	if shutdownHandler != nil {
		select {
		case shutdownHandler <- shutdownSignals[0]:
			return true
		default:
		}
	}
	return false
}
