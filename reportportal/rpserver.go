package reportportal

import (
	"github.com/avarabyeu/goRP/conf"
	"github.com/avarabyeu/goRP/registry"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

//RpServer represents ReportPortal micro-service instance
type RpServer struct {
	router *gin.Engine
	conf   *conf.RpConfig
	sd     registry.ServiceDiscovery
}

//New creates new instance of RpServer struct
func New(conf *conf.RpConfig) *RpServer {
	gin.SetMode(gin.ReleaseMode)
	rp := &RpServer{
		router: gin.Default(),
		conf:   conf,
		sd:     registry.NewConsul(conf),
	}

	rp.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]string{"status": "UP"})
	})
	return rp
}

//AddRoute gives access to GIN router to add route and perform other modifications
func (rp *RpServer) AddRoute(f func(router *gin.Engine)) {
	f(rp.router)
}

//StartServer starts HTTP server
func (rp *RpServer) StartServer() {

	// listen and server on mentioned port
	registry.Register(rp.sd)
	log.Fatal(rp.router.Run(":" + strconv.Itoa(rp.conf.Server.Port)))
}
