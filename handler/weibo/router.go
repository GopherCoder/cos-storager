package weibo

import "github.com/gin-gonic/gin"

func Register(r *gin.RouterGroup) {

	r.POST("/upload")
	r.GET("/images/:id")
	r.GET("/images_all")
	r.DELETE("/images/:id")
}

// upload images api
//"http://picupload.service.weibo.com/interface/pic_upload.php?ori=1&mime=image%2Fjpeg&data=base64&url=0&markpos=1&logo=&nick=0&marks=1&app=miniblog"
