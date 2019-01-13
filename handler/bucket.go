package handler

import (
	"cos-storager/model"
	"cos-storager/pkg/database"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func Register(r *gin.RouterGroup) {
	r.POST("/buckets", PostBucketsHandler)
	r.GET("/buckets_all", GetBucketsAllHandler)
	r.DELETE("/buckets/:id", DeleteBucketHandler)
}

type PostBucketPrams struct {
	Name string `form:"name" json:"name" binding:"required"`
	Type string `form:"type" json:"type" binding:"required"`
	Link string `form:"link" json:"link" binding:"required"`
}

func PostBucketsHandler(context *gin.Context) {
	var params PostBucketPrams
	if err := context.ShouldBind(&params); err != nil {
		return
	}
	var one model.Bucket
	if dbError := database.POSTGRES.Where("bucket_type = ?", params.Type).First(&one).Error; dbError == nil {
		MakeResponseJsonFail(context, http.StatusBadRequest, errors.New("type exists already").Error())
		return
	}
	one = model.Bucket{
		BucketName: params.Name,
		BucketType: params.Name,
		BucketURL:  params.Link,
	}
	database.POSTGRES.Save(&one)
	MakeResponseJson(context, http.StatusOK, one.BasicSerializer())
}

func GetBucketsAllHandler(context *gin.Context) {
	var records model.Buckets

	if dbError := database.POSTGRES.Model(&records).Find(&records).Error; dbError != nil {
		MakeResponseJsonFail(context, http.StatusBadRequest, dbError.Error())
		return
	}
	var results []model.BucketMessage
	for _, record := range records {
		results = append(results, record.BasicSerializer())
	}
	MakeResponseJson(context, http.StatusOK, results)
}

func DeleteBucketHandler(context *gin.Context) {
	var params string
	params = context.Param("id")
	var one model.Bucket
	if dbError := database.POSTGRES.Where("id = ?", params).First(&one).Error; dbError != nil {
		MakeResponseJsonFail(context, http.StatusBadRequest, dbError.Error())
		return
	}
	database.POSTGRES.Delete(&one)
	MakeResponseJson(context, http.StatusOK, one.BasicSerializer())
}
