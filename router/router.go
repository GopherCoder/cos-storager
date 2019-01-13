package router

import (
	"cos-storager/handler"
	"cos-storager/handler/qiniuyun"
	"cos-storager/handler/smms"
	"cos-storager/handler/tenxunyun"
	"cos-storager/handler/upyun"
	"cos-storager/handler/weibo"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Routers struct {
}

func (r *Routers) Load(g *gin.Engine, handlers ...gin.HandlerFunc) *gin.Engine {
	g.Use(gin.Recovery())
	g.Use(handlers...)

	g.NoRoute(func(context *gin.Context) {
		context.JSON(http.StatusBadRequest, gin.H{
			"data": "no exists router",
		})
	})
	fmt.Println(os.Getwd())
	g.Static("public/images", "./public")

	// Buckets
	bucketGroup := g.Group("/v1/api")
	{
		handler.Register(bucketGroup)
	}
	// 腾讯云存储
	txGroup := g.Group("/v1/api/tenxunyun")

	{
		tenxunyun.Register(txGroup)
	}

	// 七牛云存储
	qnGroup := g.Group("/v1/api/qiniuyun")
	{
		qiniuyun.Register(qnGroup)
	}

	// 又拍云
	uyGroup := g.Group("/v1/api/upyun")
	{
		upyun.Register(uyGroup)
	}

	// smms
	smGroup := g.Group("/v1/api/smms")
	{
		smms.Register(smGroup)
	}

	// 微博存储
	wbGroup := g.Group("/v1/api/weibo")
	{
		weibo.Register(wbGroup)
	}
	return g

}
