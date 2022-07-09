package webutil

import (
	"github.com/gin-gonic/gin"
	"log"
	"openid/config"
	"openid/process/queueutil"
	"os"
)

func Init() {
	log.SetOutput(os.Stdout)

	// if debug
	if !config.Server.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	// init queue
	queueutil.Init()
	// init gin
	r := gin.New()
	initRoute(r)

	log.Printf("Server running at %s ", config.Server.Addr)
	if err := r.Run(config.Server.Addr); err != nil {
		log.Panic(err)
	}
}
