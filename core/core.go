package core

import (
	"log"

	"github.com/soxft/openid-go/process/dbutil"
	"github.com/soxft/openid-go/process/queueutil"
	"github.com/soxft/openid-go/process/redisutil"
	"github.com/soxft/openid-go/process/webutil"
)

func Init() {
	log.Printf("Server initailizing...")

	// init redis
	redisutil.Init()

	// init db
	dbutil.Init()

	// init queue
	queueutil.Init()

	// init web
	webutil.Init()
}
