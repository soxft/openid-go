package main

import (
	"github.com/soxft/openid-go/process/dbutil"
	"github.com/soxft/openid-go/process/queueutil"
	"github.com/soxft/openid-go/process/redisutil"
	"github.com/soxft/openid-go/process/webutil"
)

func main() {
	redisutil.Init()
	dbutil.Init()
	queueutil.Init()
	webutil.Init()
}
