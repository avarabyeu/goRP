package server

import (
	"github.com/avarabyeu/goRP/conf"
	"github.com/avarabyeu/goRP/registry"
	"goji.io"
	"goji.io/pat"

	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

var jsonContentTypeValue = []string{"application/json; charset=utf-8"}

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
		sd = registry.NewConsul(cfg)
	}

	srv := &RpServer{
		mux: goji.NewMux(),
		cfg: cfg,
		Sd:  sd,
	}

	srv.mux.HandleFunc(pat.Get("/health"), func(w http.ResponseWriter, rq *http.Request) {
		WriteJSON(w, 200, map[string]string{"status": "UP"})
	})
	srv.mux.HandleFunc(pat.Get("/info"), func(w http.ResponseWriter, rq *http.Request) {
		WriteJSON(w, 200, map[string]interface{}{"build": map[string]string{"name": cfg.AppName}})

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

//WriteJSON serializes body to provided writer
func WriteJSON(w http.ResponseWriter, status int, body interface{}) error {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = jsonContentTypeValue
	}
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(body)
}
