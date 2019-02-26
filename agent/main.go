package main

import (
	"context"
	"fmt"
	"github.com/avarabyeu/goRP/agent/handlers"
	"github.com/caarlos0/env"
	log "github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"net/http"
)

type conf struct {
	Port     int    `env:"PORT" envDefault:"9999"`
	BasePath string `env:"BASE_PATH" envDefault:"/api/v1"`
}

func main() {
	app := fx.New(
		fx.Provide(
			newConf,
			newMux,
		),
		fx.Invoke(initServer),
	)

	app.Run()
}

func newConf() (*conf, error) {
	var cfg conf
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func newMux(cfg *conf) http.Handler {
	//TODO to be done
	return handlers.NewMux(cfg.BasePath, nil, nil)
}

func initServer(lc fx.Lifecycle, handler http.Handler, cfg *conf) {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: handler,
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Infof("Starting HTTP server on port %d", cfg.Port)

			go func() {
				if err := server.ListenAndServe(); nil != err {
					log.Fatal(err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("Stopping HTTP server.")
			return server.Shutdown(ctx)
		},
	})
}
