package main

import (
	"net/http"
	"runtime/debug"

	log "github.com/sirupsen/logrus"

	"github.com/scratch-net/ethproxy/config"
	"github.com/scratch-net/ethproxy/service"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("panic: %s: %s", r, string(debug.Stack()))
		}
	}()

	if err := Run(); err != nil && err != http.ErrServerClosed {
		log.Errorf("error while running service: %+v", err)
	}
}

func Run() error {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	logLevel, err := log.ParseLevel(cfg.LogLevel)
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(logLevel)

	srvc, err := service.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	return srvc.Run()
}
