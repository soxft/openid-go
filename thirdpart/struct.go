package thirdpart

import "github.com/gin-gonic/gin"

type ThirdPart interface {
	Handler(c *gin.Context) (string, error)
}

type ThirdProvider string

const (
	ThirdGithub ThirdProvider = "github"
)
