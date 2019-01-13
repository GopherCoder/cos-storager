package tenxunyun

import (
	"cos-storager/handler"
	"errors"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

func PostHandler(context *gin.Context) {
	var params PostParams
	if err := context.ShouldBind(&params); err != nil {
		return
	}
	files, err := context.FormFile("file_path")
	if files == nil && err != nil {
		handler.MakeResponseJsonFail(context, http.StatusBadRequest, errors.New("you should add file").Error())
		return
	}

	fileNameToStorageTenXun := handler.MakeMd5(filepath.Base(files.Filename)) + filepath.Base(files.Filename)
	var localPath string
	{
		localPath = "./public/images/tenxun/" + fileNameToStorageTenXun
	}
	if err := context.SaveUploadedFile(files, localPath); err != nil {
		handler.MakeResponseJsonFail(context, http.StatusBadRequest, err.Error())
		return
	}

	var tx TenXun
	ok, results := tx.Upload(localPath)
	if !ok {
		handler.MakeResponseJsonFail(context, http.StatusBadRequest, errors.New("upload file fail").Error())
		return
	}
	handler.MakeResponseJson(context, http.StatusOK, results)

}

func GetHandler(context *gin.Context) {
	var params int
	params, _ = strconv.Atoi(context.Param("id"))
	var (
		tx      TenXun
		ok      bool
		results interface{}
	)
	if ok, results = tx.GetOne(uint(params)); !ok {
		handler.MakeResponseJsonFail(context, http.StatusBadRequest, errors.New("get file fail").Error())
		return
	}
	handler.MakeResponseJson(context, http.StatusOK, results)
}

func GetAllHandler(context *gin.Context) {

	var (
		tx      TenXun
		ok      bool
		results interface{}
	)
	if ok, results = tx.GetAll(); !ok {
		handler.MakeResponseJsonFail(context, http.StatusBadRequest, errors.New("get all files fail").Error())
		return
	}
	handler.MakeResponseJson(context, http.StatusOK, results)

}

func DeleteOneHandler(context *gin.Context) {
	var params int
	params, _ = strconv.Atoi(context.Param("id"))

	var (
		tx      TenXun
		ok      bool
		results interface{}
	)
	if ok, results = tx.DeleteOne(uint(params)); !ok {
		handler.MakeResponseJsonFail(context, http.StatusBadRequest, errors.New("delete file fail").Error())
		return
	}
	handler.MakeResponseJson(context, http.StatusOK, results)

}

func DeleteAllHandler(context *gin.Context) {

	var (
		tx      TenXun
		ok      bool
		results interface{}
	)
	if ok, results = tx.DeleteAll(); !ok {
		handler.MakeResponseJsonFail(context, http.StatusBadRequest, errors.New("delete all files fail").Error())
		return
	}
	handler.MakeResponseJson(context, http.StatusOK, results)
}
