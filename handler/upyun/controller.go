package upyun

import (
	"cos-storager/handler"
	"cos-storager/model"
	"cos-storager/pkg/database"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetUpYunByIDHandler(context *gin.Context) {
	var params GetUpYunOneParams
	if err := context.ShouldBindQuery(&params); err != nil {
		return
	}
	if params.ID == 0 && params.Path == "" {
		handler.MakeResponseJsonFail(context, http.StatusBadRequest, errors.New("add id or path in path params").Error())
		return
	}
	var upyun UPYun
	ok, results := upyun.GetOne(params.ID, params.Path)
	if !ok {
		handler.MakeResponseJsonFail(context, http.StatusBadRequest, errors.New("get images info by upyun fail"))
		return
	}
	handler.MakeResponseJson(context, http.StatusOK, results)
}

func GetUpYunByPathHandler(context *gin.Context) {
	var params GetUpYunByPathParams
	if err := context.ShouldBindQuery(&params); err != nil {
		return
	}
	if params.Path == "" {
		params.Path = "/images/"
	}
	var filesMessages model.FilesMessages
	if dbError := database.POSTGRES.Where("files_message_name ILIKE ?", fmt.Sprintf("%%%s%%", params.Path)).Find(&filesMessages).Error; dbError != nil {
		handler.MakeResponseJsonFail(context, http.StatusBadRequest, errors.New("can not find records").Error())
		return
	}
	var results model.Messages
	for _, one := range filesMessages {
		results = append(results, one.BasicSerialize())
	}
	handler.MakeResponseJson(context, http.StatusOK, results)

}

func GetListStorageHandler(context *gin.Context) {
	var upy UPYun
	_, results := upy.ListPath()
	handler.MakeResponseJson(context, http.StatusOK, results)
}

func GetUpYunAllHandler(context *gin.Context) {
	var bucket model.Bucket
	if dbError := database.POSTGRES.Where("bucket_type = ?", "upyun").First(&bucket).Error; dbError != nil {
		handler.MakeResponseJsonFail(context, http.StatusBadRequest, errors.New("not exists this bucket").Error())
		return
	}
	var fileMessages model.FilesMessages
	if dbError := database.POSTGRES.Where("bucket_id = ?", bucket.ID).Find(&fileMessages).Error; dbError != nil {
		handler.MakeResponseJson(context, http.StatusOK, fileMessages)
		return
	}
	var results model.Messages
	for _, one := range fileMessages {
		results = append(results, one.BasicSerialize())
	}
	handler.MakeResponseJson(context, http.StatusOK, results)
}

func PostUpYunHandler(context *gin.Context) {
	var params PostUpYunParams
	if err := context.ShouldBind(&params); err != nil {
		return
	}
	files, err := context.FormFile("file_path")
	if files == nil && err != nil {
		handler.MakeResponseJsonFail(context, http.StatusBadRequest, errors.New("you should add file").Error())
		return
	}

	fileNameToStorageTenXun := handler.MakeMd5(filepath.Base(files.Filename))
	typeList := strings.Split(files.Filename, ".")
	imageType := typeList[len(typeList)-1]
	var localPath string
	{
		localPath = "./public/images/upyun/" + fileNameToStorageTenXun + "." + imageType
	}
	if err := context.SaveUploadedFile(files, localPath); err != nil {
		handler.MakeResponseJsonFail(context, http.StatusBadRequest, err.Error())
		return
	}

	if params.Path == "" {
		params.Path = "/images/" + fileNameToStorageTenXun + "." + imageType
	} else {
		params.Path = fmt.Sprintf("/%s/%s", params.Path, fileNameToStorageTenXun+"."+imageType)
	}
	var upyun UPYun
	ok, results := upyun.Upload(params.Path, localPath)
	if !ok {
		handler.MakeResponseJsonFail(context, http.StatusBadRequest, errors.New("upload images by upyun fail").Error())
		return
	}
	os.Remove(localPath)
	handler.MakeResponseJson(context, http.StatusOK, results)
}

func GetUsageHandler(context *gin.Context) {
	var upy UPYun
	results, err := upy.Usage()
	if err != nil {
		handler.MakeResponseJsonFail(context, http.StatusBadRequest, err.Error())
		return
	}
	handler.MakeResponseJson(context, http.StatusOK, results)

}

func DeleteUpYunHandler(context *gin.Context) {
	var params DeleteUpYunParams
	if err := context.ShouldBindQuery(&params); err != nil {
		return
	}
	if params.ID == 0 && params.Path == "" {
		handler.MakeResponseJsonFail(context, http.StatusBadRequest, errors.New("id or path should not be nil").Error())
		return
	}
	var upy UPYun
	ok, results := upy.DeleteImages(params.ID, params.Path)
	if !ok {
		handler.MakeResponseJsonFail(context, http.StatusBadRequest, errors.New("record delete fail or not exists this record").Error())
		return
	}
	handler.MakeResponseJson(context, http.StatusOK, results)

}
