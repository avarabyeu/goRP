package server

import (
	"github.com/avarabyeu/goRP/conf"
	"goji.io"
	"goji.io/pat"
	"net/http"
	"os"
	"github.com/gorilla/handlers"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func ExampleRpServer() {
	rpConf := conf.LoadConfig("../server.yaml", nil)
	rp := New(rpConf)

	rp.AddRoute(func(router *goji.Mux) {
		router.HandleFunc(pat.Get("/ping"), func(w http.ResponseWriter, rq *http.Request) {
			WriteJSON(w, http.StatusOK, Person{"av", 20})
		})
	})

	rp.StartServer()

}

func ExampleAuthMiddleware() {

	rpConf := conf.LoadConfig("../server.yaml",
		map[string]interface{}{"AuthServerURL": "http://localhost:9998/sso/me"})

	srv := New(rpConf)

	srv.AddRoute(func(mux *goji.Mux) {
		mux.Use(func(next http.Handler) http.Handler {
			return handlers.LoggingHandler(os.Stdout, next)
		})

		secured := goji.SubMux()
		secured.Use(RequireRole("USER", rpConf.Get("AuthServerURL").(string)))

		me := func(w http.ResponseWriter, rq *http.Request) {
			WriteJSON(w, http.StatusOK, rq.Context().Value("user"))

		}
		secured.HandleFunc(pat.Get("/me"), me)

		mux.Handle(pat.Get("/"), secured)

	})

	srv.StartServer()
}
