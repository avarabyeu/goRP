package reportportal

import (
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"time"
	"net/http"
)

const (
	PROTOCOL = "http://"
	RETRY_TIMEOUT time.Duration = time.Second * 5
	POLL_INTERVAL time.Duration = time.Second * 15
	RETRY_ATTEMPTS int = 3
)

type RpServer struct {
	router *gin.Engine
	conf   *RpConfig
}

func New(conf *RpConfig) *RpServer {
	gin.SetMode(gin.ReleaseMode)
	rp := &RpServer{
		router: gin.Default(),
		conf: conf,
	}

	rp.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]string{"status": "UP" })
	})
	return rp
}

func (rp *RpServer) AddRoute(f func(router *gin.Engine)) {
	f(rp.router)
}

func (rp *RpServer) StartServer() {

	// listen and server on mentioned port
	registerInEureka(rp)

	log.Fatal(rp.router.Run(":" + strconv.Itoa(rp.conf.Server.Port)))
}


