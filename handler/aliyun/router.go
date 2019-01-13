package aliyun

import "github.com/gin-gonic/gin"

func Register(r *gin.RouterGroup) {
	r.GET("/images", GetAliYunByIDHandler)
	r.GET("/images_all", GetAliYunAllHandler)
	r.DELETE("/images/:id", DeleteAliYunByIDHandler)
	r.POST("/images", PostAliYunUploadHandler)
	r.GET("/images_usage", GetAliYunUsageHandler)
	r.GET("/images_storage", GetAliYunStorageHandler)
}
