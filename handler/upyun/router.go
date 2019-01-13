package upyun

import "github.com/gin-gonic/gin"

func Register(r *gin.RouterGroup) {
	r.GET("/images", GetUpYunByIDHandler)
	r.GET("/images_path", GetUpYunByPathHandler)
	r.GET("/images_storage", GetListStorageHandler)
	r.GET("/images_all", GetUpYunAllHandler)
	r.POST("/upload", PostUpYunHandler)
	r.GET("/images_usage", GetUsageHandler)
	r.DELETE("/images_remove", DeleteUpYunHandler)
}
