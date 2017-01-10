package main

import (
	"github.com/avarabyeu/goRP/conf"
	"github.com/avarabyeu/goRP/reportportal"
	"github.com/gin-gonic/gin"
	"net/http"
)

type person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {

	rpConf := conf.LoadConfig("server.yaml")
	rp := reportportal.New(rpConf)

	rp.AddRoute(func(router *gin.Engine) {
		router.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, person{"av", 20})
		})
	})

	rp.StartServer()

}
