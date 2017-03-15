package server

import (
	"github.com/avarabyeu/goRP/conf"
	"github.com/avarabyeu/goRP/registry"
	"goji.io"
	"goji.io/pat"

	"log"
	"net/http"
	"strconv"
	"github.com/avarabyeu/goRP/commons"
)

//RpServer represents ReportPortal micro-service instance
type RpServer struct {
	mux *goji.Mux
	cfg *conf.RpConfig
	Sd  registry.ServiceDiscovery
}

//New creates new instance of RpServer struct
func New(cfg *conf.RpConfig) *RpServer {

	var sd registry.ServiceDiscovery
	switch cfg.Registry {
	case conf.Eureka:
		sd = registry.NewEureka(cfg)
	case conf.Consul:
		cfg.Consul.Tags = cfg.Consul.Tags + " statusPageUrlPath=/info" + " " + "healthCheckUrlPath=/health"
		sd = registry.NewConsul(cfg)
	}

	srv := &RpServer{
		mux: goji.NewMux(),
		cfg: cfg,
		Sd:  sd,
	}

	srv.mux.HandleFunc(pat.Get("/health"), func(w http.ResponseWriter, rq *http.Request) {
		commons.WriteJSON(200, map[string]string{"status": "UP"}, w)
	})
	srv.mux.HandleFunc(pat.Get("/info"), func(w http.ResponseWriter, rq *http.Request) {
		commons.WriteJSON(200, map[string]interface{}{"build": map[string]string{"name": cfg.AppName}}, w)

	})
	return srv
}

//AddRoute gives access to GIN router to add route and perform other modifications
func (srv *RpServer) AddRoute(f func(router *goji.Mux)) {
	f(srv.mux)
}

//StartServer starts HTTP server
func (srv *RpServer) StartServer() {

	if nil != srv.Sd {
		registry.Register(srv.Sd)
	}
	// listen and server on mentioned port
	log.Printf("Starting on port %d", srv.cfg.Server.Port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(srv.cfg.Server.Port), srv.mux))
}
