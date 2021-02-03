package rest

import (
	"github.com/gin-gonic/gin"

	"github/demo/rest/device"
)

func Init() *gin.Engine {
	r := gin.New()

	r.Use(gin.Recovery())

	v1 := r.Group("/v1")
	{
		device.MakeHandler(v1)
	}

	return r
}
