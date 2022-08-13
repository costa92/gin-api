package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/costa92/go-web/cmd"
	"github.com/costa92/go-web/config"
	"github.com/costa92/go-web/internal/db"
	"github.com/costa92/go-web/server"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Info().Msg("start service")

	cmd.Execute()
	cfg, _ := config.NewConfig()
	// binding.Validator = new(validator.DefaultValidator)
	db.InitDB(cfg)
	app := server.NewServer(cfg.ServerConf)
	err := app.Run()
	if err != nil {
		panic(err)
	}
}
