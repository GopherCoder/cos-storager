package smms

import (
	"cos-storager/config"
	"cos-storager/handler"
	"cos-storager/model"
	"cos-storager/pkg/database"
	"errors"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func PostSMMSHandler(context *gin.Context) {

	smFiles, err := context.FormFile("smfile")
	if smFiles == nil && err != nil {
		handler.MakeResponseJsonFail(context, http.StatusBadRequest, errors.New("you should add file").Error())
		return
	}
	fileNameToStorageTenXun := handler.MakeMd5(filepath.Base(smFiles.Filename))
	typeList := strings.Split(smFiles.Filename, ".")
	imageType := typeList[len(typeList)-1]
	var localPath string
	{
		localPath = "./public/images/smms/" + fileNameToStorageTenXun + "." + imageType
	}
	if err := context.SaveUploadedFile(smFiles, localPath); err != nil {
		handler.MakeResponseJsonFail(context, http.StatusBadRequest, err.Error())
		return
	}
	smms := NewSmmS()
	ok, result := smms.Upload(localPath)
	if !ok {
		handler.MakeResponseJsonFail(context, http.StatusBadRequest, errors.New("upload file fail").Error())
		return
	}
	handler.MakeResponseJson(context, http.StatusOK, result)

}

func GetSMMSByIDHandler(context *gin.Context) {
	var params GetSMMSImageParams
	if err := context.ShouldBindQuery(&params); err != nil {
		return
	}

	var bucket model.Bucket
	if dbError := database.POSTGRES.Where("bucket_type = ?", "smms").First(&bucket).Error; dbError != nil {
		handler.MakeResponseJsonFail(context, http.StatusBadRequest, dbError.Error())
		return
	}

	var fileMessage model.FilesMessage
	if dbError := database.POSTGRES.Where("id = ? AND bucket_id = ?", params.ID, bucket.ID).First(&fileMessage).Error; dbError != nil {
		handler.MakeResponseJsonFail(context, http.StatusBadRequest, dbError.Error())
		return
	}

	handler.MakeResponseJson(context, http.StatusOK, fileMessage.BasicSerialize())

}

func GetSMMSAllHandler(context *gin.Context) {

	var bucket model.Bucket
	if dbError := database.POSTGRES.Where("bucket_type = ?", config.SMMSBUCKET).First(&bucket).Error; dbError != nil {
		handler.MakeResponseJsonFail(context, http.StatusBadRequest, dbError.Error())
		return
	}

	var fileMessages model.FilesMessages
	if dbError := database.POSTGRES.Where("bucket_id = ?", bucket.ID).Find(&fileMessages).Error; dbError != nil {
		handler.MakeResponseJsonFail(context, http.StatusBadRequest, dbError.Error())
		return
	}
	var results model.Messages
	for _, i := range fileMessages {
		results = append(results, i.BasicSerialize())
	}
	handler.MakeResponseJson(context, http.StatusOK, results)

}

func DeleteSMMSHandler(context *gin.Context) {
	var params string
	params = context.Param("id")
	var fileMessage model.FilesMessage
	if dbError := database.POSTGRES.Where("id = ?", params).First(&fileMessage).Error; dbError != nil {
		handler.MakeResponseJsonFail(context, http.StatusBadRequest, dbError.Error())
		return
	}
	database.POSTGRES.Delete(&fileMessage)
	var results = make(map[string]interface{})
	results["fileMessages"] = fileMessage
	results["fileInfos"] = fileMessage.BasicSerialize()

	handler.MakeResponseJson(context, http.StatusOK, results)

}
