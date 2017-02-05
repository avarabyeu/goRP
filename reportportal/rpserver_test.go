package reportportal

import (
	"github.com/avarabyeu/goRP/conf"
	"net/http"
	"goji.io"
	"goji.io/pat"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func ExampleRpServer() {
	rpConf := conf.LoadConfig("../server.yaml")
	rp := New(rpConf)

	rp.AddRoute(func(router *goji.Mux) {
		router.HandleFunc(pat.Get("/ping"), func(w http.ResponseWriter, rq *http.Request) {
			WriteJSON(w, http.StatusOK, Person{"av", 20})
		})
	})

	rp.StartServer()

}
