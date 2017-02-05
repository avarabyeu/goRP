package reportportal

import (
	"github.com/avarabyeu/goRP/conf"
	"github.com/avarabyeu/goRP/registry"
	"goji.io"
	"goji.io/pat"

	"net/http"
	"encoding/json"
	"log"
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
	rp := &RpServer{
		mux:    goji.NewMux(),
		conf:   conf,
		sd:     registry.NewConsul(conf),
	}

	rp.mux.HandleFunc(pat.Get("/health"), func(w http.ResponseWriter, rq *http.Request) {
		WriteJSON(w, 200, map[string]string{"status": "UP"})
	})
	return rp
}

//AddRoute gives access to GIN router to add route and perform other modifications
func (rp *RpServer) AddRoute(f func(router *goji.Mux)) {
	f(rp.mux)
}

//StartServer starts HTTP server
func (rp *RpServer) StartServer() {
	// listen and server on mentioned port
	//registry.Register(rp.sd)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(rp.conf.Server.Port), rp.mux))
}

func WriteJSON(w http.ResponseWriter, status int, body interface{}) error {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = jsonContentTypeValue
	}
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(body)
}
