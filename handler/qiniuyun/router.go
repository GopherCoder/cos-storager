package qiniuyun

import "github.com/gin-gonic/gin"

func Register(r *gin.RouterGroup) {

	r.POST("/upload", PostUploadHandler)
	r.GET("/images/:id", GetImageHandler)
	r.GET("/images_all", GetImagesHandler)
	r.DELETE("/images/:id", DeleteOneImageHandler)
	r.DELETE("/images_remove", DeleteImagesHandler)
}
