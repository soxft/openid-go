package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"openid/config"
	"openid/route"
	"os"
)

func main() {
	// if debug
	if !config.C.Server.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	// init gin
	r := gin.New()
	route.Init(r)

	log.Printf("Server running at %s ", config.C.Server.Addr)
	if err := r.Run(config.C.Server.Addr); err != nil {
		log.Panic(err)
	}
}

func init() {
	log.SetOutput(os.Stdout)
}
