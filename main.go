package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"openid/config"
	"openid/process/queueutil"
	"openid/route"
	"os"
)

func main() {
	log.SetOutput(os.Stdout)

	// if debug
	if !config.C.Server.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	// init queue
	queueutil.Init()
	// init gin
	r := gin.New()
	route.Init(r)

	log.Printf("Server running at %s ", config.C.Server.Addr)
	if err := r.Run(config.C.Server.Addr); err != nil {
		log.Panic(err)
	}
}
