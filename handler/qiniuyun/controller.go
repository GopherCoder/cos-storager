package qiniuyun

import (
	"cos-storager/handler"
	"cos-storager/model"
	"cos-storager/pkg/database"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

func PostUploadHandler(context *gin.Context) {
	var params PostUploadParams
	if err := context.ShouldBind(&params); err != nil {
		return
	}

	files, err := context.FormFile("file_path")
	if files == nil && err != nil {
		handler.MakeResponseJsonFail(context, http.StatusBadRequest, errors.New("you should add file").Error())
		return
	}

	fileNameToStorageQiniu := handler.MakeMd5(filepath.Base(files.Filename)) + filepath.Base(files.Filename)
	var localPath string
	{
		localPath = "./public/images/qiniu/" + fileNameToStorageQiniu
	}
	if err := context.SaveUploadedFile(files, localPath); err != nil {
		handler.MakeResponseJsonFail(context, http.StatusBadRequest, err.Error())
		return
	}

	var qiniu QiNiu
	ok, result := qiniu.Upload(localPath)
	if err := os.Remove(localPath); err != nil {
		handler.MakeResponseJsonFail(context, http.StatusBadRequest, err.Error())
		return
	}
	if !ok {
		handler.MakeResponseJsonFail(context, http.StatusBadRequest, errors.New("upload file fail").Error())
		return
	}
	resultFromQiNiu := result.(QiNiuResult)

	var oneBucket model.Bucket
	if dbError := database.POSTGRES.Where("bucket_type = ?", resultFromQiNiu.Bucket).First(&oneBucket).Error; dbError != nil {
		handler.MakeResponseJsonFail(context, http.StatusBadRequest, dbError.Error())
		return
	}
	var one model.FilesMessage
	one = model.FilesMessage{
		FilesMessageName: resultFromQiNiu.Name,
		FilesMessageKey:  resultFromQiNiu.Key,
		FilesMessageSize: resultFromQiNiu.Fsize,
		FilesMessageURL:  fmt.Sprintf(oneBucket.BucketURL+"/%s", resultFromQiNiu.Key),
		BucketID:         oneBucket.ID,
	}
	database.POSTGRES.Save(&one)
	handler.MakeResponseJson(context, http.StatusOK, one.BasicSerialize())
}

func GetImageHandler(context *gin.Context) {
	var params int
	params, _ = strconv.Atoi(context.Param("id"))
	var qiniu QiNiu
	ok, result := qiniu.GetOne(uint(params))
	if !ok {
		handler.MakeResponseJsonFail(context, http.StatusBadRequest, errors.New("can not find this record").Error())
		return
	}
	resultInfo := result.(map[string]interface{})

	handler.MakeResponseJson(context, http.StatusOK, resultInfo)
}

func GetImagesHandler(context *gin.Context) {

	var qiniu QiNiu
	ok, results := qiniu.GetAll()
	if !ok {
		handler.MakeResponseJsonFail(context, http.StatusBadRequest, errors.New("can not found records").Error())
		return
	}
	allResults := results.(map[string]interface{})
	var mes model.Messages
	for _, one := range allResults["fileMessages"].(model.FilesMessages) {
		mes = append(mes, one.BasicSerialize())
	}
	allResults["urls"] = mes
	handler.MakeResponseJson(context, http.StatusOK, allResults)
}

func DeleteOneImageHandler(context *gin.Context) {
	var params int
	params, _ = strconv.Atoi(context.Param("id"))

	var qiniu QiNiu
	ok, results := qiniu.DeleteOne(uint(params))
	if !ok {
		handler.MakeResponseJsonFail(context, http.StatusBadRequest, errors.New("can not delete record").Error())
		return
	}
	handler.MakeResponseJson(context, http.StatusOK, results)
}

func DeleteImagesHandler(context *gin.Context) {
	var qiniu QiNiu
	ok, results := qiniu.DeleteAll()
	if !ok {
		handler.MakeResponseJsonFail(context, http.StatusBadRequest, errors.New("can not delete all files").Error())
		return
	}
	handler.MakeResponseJson(context, http.StatusOK, results)
}
