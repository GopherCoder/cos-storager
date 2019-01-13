package smms

import "github.com/gin-gonic/gin"

func Register(r *gin.RouterGroup) {
	r.POST("/upload", PostSMMSHandler)
	r.GET("/images", GetSMMSByIDHandler)
	r.GET("/images_all", GetSMMSAllHandler)
	r.DELETE("/images/:id", DeleteSMMSHandler)
}
