package device

import "github.com/gin-gonic/gin"

func MakeHandler(r *gin.RouterGroup) {
	g := r.Group("/device")
	{
		g.GET("", FindDevice)
		g.POST("", RegisterDevice)
		g.DELETE("/:id", DeleteDevice)
		g.PUT("", UpdateDevice)
	}
}
