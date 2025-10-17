package app

import (
	"os"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var (
	router *gin.Engine
)

func init() {
	router = gin.Default()
}

func StartRoute() {
	MapMiddewares()
	MapURLs()

	host := os.Getenv("APP_HOST")
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	hostport := host + ":" + port

	log.Info("Iniciando servidor en: ", hostport)
	router.Run(hostport)
}
