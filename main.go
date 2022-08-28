package main

import (
	"github.com/costa92/go-web/cmd"
	"github.com/costa92/go-web/config"
	"github.com/costa92/go-web/internal/db"
	"github.com/costa92/go-web/internal/logger"
	"github.com/costa92/go-web/server"
)

func main() {
	cmd.Execute()
	cfg, _ := config.NewConfig()
	// 初始化日志
	logger.Init(cfg.Logger)
	defer logger.Flush()
	// 初始化数据库
	db.InitDB(cfg)
	app := server.NewServer(cfg.ServerConf)
	err := app.Run()
	if err != nil {
		panic(err)
	}
}
