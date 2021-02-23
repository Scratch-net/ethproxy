package service

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	"golang.org/x/sync/errgroup"

	"github.com/scratch-net/ethproxy/cache"
	"github.com/scratch-net/ethproxy/config"
	"github.com/scratch-net/ethproxy/controller"
	"github.com/scratch-net/ethproxy/fetcher"
)

type Service struct {
	server *http.Server
}

func New(cfg *config.Config) (*Service, error) {
	log.Info("starting service...")

	cli := fetcher.NewCloudflareEthClient(cfg.ClientTimeout)

	cacheStorage, err := cache.NewBlockCacheStorage(cli, cfg.CacheSize)
	if err != nil {
		return nil, err
	}

	// load latest block to make sure everything works
	block, err := cacheStorage.Get(cache.LatestBlock)
	if err != nil {
		return nil, err
	}
	if block == nil {
		return nil, errors.New("cache returned nil latest block")
	}

	ctrl, err := controller.New(cacheStorage)
	if err != nil {
		return nil, err
	}
	handler, err := NewHandler(ctrl)
	if err != nil {
		return nil, err
	}
	router, err := NewRouter(handler)
	if err != nil {
		return nil, err
	}

	srv := &http.Server{
		Handler:      router,
		Addr:         cfg.ListenAddress,
		WriteTimeout: cfg.WriteTimeout,
		ReadTimeout:  cfg.ReadTimeout,
	}

	return &Service{
		server: srv,
	}, nil
}

func (s *Service) Run() (err error) {

	errGroup, ctx := errgroup.WithContext(context.Background())
	errGroup.Go(func() error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case sig := <-gracefulShutdown():
			log.Infof("got signal: %s", sig)
			return s.server.Shutdown(ctx)
		}
	})

	errGroup.Go(func() error {
		return s.server.ListenAndServe()
	})

	log.Info("service started")

	return errGroup.Wait()

}

func gracefulShutdown() chan os.Signal {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	return signalChannel
}
