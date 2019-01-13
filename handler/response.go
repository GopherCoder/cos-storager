package handler

import (
	"crypto/md5"
	"encoding/hex"
	"time"

	"github.com/gin-gonic/gin"
)

type OperationCos interface {
	Upload(filePath string) (bool, interface{})
	GetOne(id uint) (bool, interface{})
	GetAll() (bool, interface{})
	DeleteOne(id uint) (bool, interface{})
	DeleteAll() (bool, interface{})
}

func MakeResponseJson(context *gin.Context, code int, values interface{}) {
	context.JSON(code, gin.H{
		"data": values,
	})
}

func MakeResponseJsonFail(context *gin.Context, code int, values interface{}) {
	context.JSON(code, gin.H{
		"err": values,
	})
}

func MakeReponseString(context *gin.Context, code int, values interface{}) {
	context.String(code, values.(string))
}

func MakeMd5(key string) string {
	hash := md5.New()
	hash.Write([]byte(key + time.Now().String()))
	return hex.EncodeToString(hash.Sum(nil))
}
