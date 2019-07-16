package main

import (
	"context"
	"fmt"
	"github.com/avarabyeu/goRP/agent/store"
	"github.com/avarabyeu/goRP/gorp"
	"github.com/dgraph-io/badger"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/avarabyeu/goRP/agent/handlers"
	"github.com/caarlos0/env"
	log "github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

type conf struct {
	Port           int    `env:"PORT" envDefault:"9999"`
	BasePath       string `env:"BASE_PATH" envDefault:"/api/v1"`
	RpProxyURL     string `env:"RP_PROXY_URL" envDefault:"http://localhost:8080"`
	RpProxyProject string `env:"RP_PROXY_PROJECT" envDefault:"default_personal"`
	RpProxyUUID    string `env:"RP_PROXY_UUID" envDefault:"01b2a730-4a96-43ae-b4c5-86db2e1338a9"`
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
	// TODO to be done
	tmpDir, err := ioutil.TempDir("", "badger")
	if err != nil {
		log.Fatal(err)
	}

	db, err := badger.Open(badger.DefaultOptions(tmpDir))
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	defer os.RemoveAll(tmpDir) // clean up

	return handlers.NewMux(cfg.BasePath, gorp.NewClient(cfg.RpProxyURL, cfg.RpProxyURL, cfg.RpProxyUUID), store.NewBadgerStore(db))
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
				if err := server.ListenAndServe(); err != nil {
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
