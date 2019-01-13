package weibo

import "github.com/gin-gonic/gin"

func Register(r *gin.RouterGroup) {

	r.POST("/upload")
	r.GET("/images/:id")
	r.GET("/images_all")
	r.DELETE("/images/:id")
}
