package v1

import "github.com/gin-gonic/gin"

type AdapterHandler struct {
	Engine *gin.Engine
	Router *gin.RouterGroup
}

func NewAdapterHandler() *AdapterHandler {
	return &AdapterHandler{
		Engine: gin.New(),
	}
}
