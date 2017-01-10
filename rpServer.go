package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/avarabyeu/gorp-commons/reportportal"
	"github.com/avarabyeu/gorp-commons/conf"
)

type Person struct {
	Name string `json:"name"`
	Age  int `json:"age"`
}

func main() {

	rpConf := conf.LoadConfig("server.yaml")
	rp := reportportal.New(rpConf)

	rp.AddRoute(func(router *gin.Engine) {
		router.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, Person{"av", 20})
		})
	})

	rp.StartServer()

}


