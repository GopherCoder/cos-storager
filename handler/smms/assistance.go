package smms

import (
	"bytes"
	"cos-storager/model"
	"cos-storager/pkg/database"
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/tidwall/gjson"
)

type SmmS struct {
	API string `json:"api"`
}

func NewSmmS() *SmmS {
	return &SmmS{
		API: "https://sm.ms/api/upload",
	}
}

func (ss SmmS) Upload(localPath string) (bool, interface{}) {

	var bucket model.Bucket
	if dbError := database.POSTGRES.Where("bucket_type = ?", "smms").First(&bucket).Error; dbError != nil {
		return false, dbError.Error()
	}

	paths := strings.Split(localPath, "/")
	name := paths[len(paths)-1]

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	fileWriter, _ := bodyWriter.CreateFormFile("smfile", name)
	file, _ := os.Open(localPath)
	defer file.Close()
	io.Copy(fileWriter, file)
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	response, _ := http.Post(ss.API, contentType, bodyBuf)
	content, _ := ioutil.ReadAll(response.Body)

	values := gjson.ParseBytes(content)
	code := values.Get("code").String()

	if strings.Contains(code, "error") {
		return false, errors.New("upload file fail")
	}
	var fileMessage model.FilesMessage
	fileMessage = model.FilesMessage{
		FilesMessageName: values.Get("data.storename").String(),
		FilesMessageURL:  values.Get("data.url").String(),
		FilesMessageSize: int(values.Get("data.size").Int()),
		FilesMessageKey:  values.Get("data.hash").String(),
		BucketID:         bucket.ID,
	}
	database.POSTGRES.Save(&fileMessage)
	os.Remove(localPath)
	var results = make(map[string]interface{})
	results["buckets"] = bucket
	results["fileMessages"] = fileMessage.BasicSerialize()
	results["fileInfo"] = fileMessage
	return true, results
}
