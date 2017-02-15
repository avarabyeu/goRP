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
	mux  *goji.Mux
	conf *conf.RpConfig
	sd   registry.ServiceDiscovery
}

//New creates new instance of RpServer struct
func New(conf *conf.RpConfig) *RpServer {
	srv := &RpServer{
		mux:  goji.NewMux(),
		conf: conf,
		sd:   registry.NewConsul(conf),
	}

	srv.mux.HandleFunc(pat.Get("/health"), func(w http.ResponseWriter, rq *http.Request) {
		WriteJSON(w, 200, map[string]string{"status": "UP"})
	})
	return srv
}

//AddRoute gives access to GIN router to add route and perform other modifications
func (srv *RpServer) AddRoute(f func(router *goji.Mux)) {
	f(srv.mux)
}

//StartServer starts HTTP server
func (srv *RpServer) StartServer() {
	// listen and server on mentioned port
	registry.Register(srv.sd)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(srv.conf.Server.Port), srv.mux))
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
