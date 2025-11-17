package dto

import "github.com/gin-gonic/gin"

type PageRequest struct {
	Page    int `json:"page" binding:"min=1"`
	PerPage int `json:"per_page" binding:"min=1,max=100"`
}

type ListResponse struct {
	Total int         `json:"total"`
	List  interface{} `json:"list"`
}

type SuccessResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type IDResponse struct {
	ID interface{} `json:"id"`
}

type DataResponse struct {
	Data interface{} `json:"data"`
}

func BindJSON(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		return err
	}
	return nil
}