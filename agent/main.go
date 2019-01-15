package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"net/http"
)

type conf struct {
	Port int
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

func newConf() *conf {
	return &conf{}
}

func newMux() http.Handler {
	mux := chi.NewMux()

	return mux
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
