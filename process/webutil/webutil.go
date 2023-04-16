package webutil

import (
	"github.com/gin-gonic/gin"
	"github.com/soxft/openid-go/config"
	"log"
	"os"
)

func Init() {
	log.SetOutput(os.Stdout)

	// if debug
	if !config.Server.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	
	// init gin
	r := gin.New()
	initRoute(r)

	log.Printf("Server running at %s ", config.Server.Addr)
	if err := r.Run(config.Server.Addr); err != nil {
		log.Panic(err)
	}
}
