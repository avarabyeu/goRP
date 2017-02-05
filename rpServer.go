package main

import (
	"github.com/avarabyeu/goRP/conf"
	"github.com/avarabyeu/goRP/reportportal"
	"net/http"
	"goji.io"
	"goji.io/pat"
	"github.com/gorilla/handlers"

	"os"
)

type person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {

	rpConf := conf.LoadConfig("server.yaml")
	rp := reportportal.New(rpConf)

	rp.AddRoute(func(mux *goji.Mux) {
		mux.Use(func(next http.Handler) http.Handler {
			return handlers.LoggingHandler(os.Stdout, next)
		})

		mux.Use(reportportal.RequireRole("USER", rpConf.AuthServerURL))

		me := func(w http.ResponseWriter, rq *http.Request) {
			reportportal.WriteJSON(w, http.StatusOK, rq.Context().Value("user"))

		}
		mux.HandleFunc(pat.Get("/me"), me)
	})

	rp.StartServer()

}
