package tenxunyun

import "github.com/gin-gonic/gin"

func Register(r *gin.RouterGroup) {

	r.POST("/upload", PostHandler)
	r.GET("/images/:id", GetHandler)
	r.GET("/images_all", GetAllHandler)
	r.DELETE("/images/:id", DeleteOneHandler)
	r.DELETE("/images_remove", DeleteAllHandler)
}
