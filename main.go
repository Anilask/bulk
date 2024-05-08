package main

import (
	"bulk/config"
	"bulk/server"

	"bulk/logger"
)

func main() {
	cfg := config.Cfg
	log := logger.New(cfg.Tracer.TracerName, cfg.ENV, cfg.LogLevel)
	sv := server.New(cfg, log)
	sv.Start()
}
