package main

import (
	"github.com/avarabyeu/goRP/conf"
	"github.com/avarabyeu/goRP/server"

	"net/http"
	"goji.io"
	"goji.io/pat"
	"github.com/gorilla/handlers"

	"os"
)

func main() {

	rpConf := conf.LoadConfig("server.yaml", nil)
	srv := server.New(rpConf)

	srv.AddRoute(func(mux *goji.Mux) {
		mux.Use(func(next http.Handler) http.Handler {
			return handlers.LoggingHandler(os.Stdout, next)
		})

		//secured := goji.SubMux()
		//secured.Use(server.RequireRole("USER", rpConf.AuthServerURL))
		//
		//
		//me := func(w http.ResponseWriter, rq *http.Request) {
		//	server.WriteJSON(w, http.StatusOK, rq.Context().Value("user"))
		//
		//}
		//secured.HandleFunc(pat.Get("/me"), me)


		dir := "/Users/avarabyeu/work/sources/reportportal/service-ui/build/resources/main/"
		//mux.Handle(pat.Get("/public/*"), http.StripPrefix("/public", http.FileServer(http.Dir(dir))))
		mux.Handle(pat.Get("/public/*"), http.FileServer(http.Dir(dir)))
		//mux.Handle(pat.Get("/"), secured)

	})

	srv.StartServer()

}
