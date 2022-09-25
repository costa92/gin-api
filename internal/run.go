package internal

import "github.com/costa92/go-web/config"

func Run(cfg *config.Config) error {
	server, err := createAPIServer(cfg)
	if err != nil {
		return err
	}
	return server.preRun().Run()
}
