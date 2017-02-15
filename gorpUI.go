package main

import (
	"github.com/avarabyeu/goRP/conf"
	"github.com/avarabyeu/goRP/server"

	"net/http"
	"goji.io"
	"goji.io/pat"
	"github.com/gorilla/handlers"

	"os"
	"log"
)

func main() {

	currDir, _ := os.Getwd()
	rpConf := conf.LoadConfig("", map[string]interface{}{"staticsPath": currDir})
	srv := server.New(rpConf)

	srv.AddRoute(func(mux *goji.Mux) {
		mux.Use(func(next http.Handler) http.Handler {
			return handlers.LoggingHandler(os.Stdout, next)
		})

		dir := rpConf.Get("staticsPath").(string)
		err := os.Chdir(dir)
		if nil != err {
			log.Fatalf("Dir %s not found", dir)
		}

		mux.Handle(pat.Get("/*"), http.FileServer(http.Dir(dir)))

	})

	srv.StartServer()

}
