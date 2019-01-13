package upyun

import (
	"cos-storager/config"
	"cos-storager/model"
	"cos-storager/pkg/database"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/upyun/go-sdk/upyun"
)

type UPYun struct {
}

func NewUpYun() *upyun.UpYun {
	return upyun.NewUpYun(&upyun.UpYunConfig{
		Bucket:   config.UPYUNBUCKET,
		Operator: config.UPYUNOPERATOR,
		Password: config.UPYUNPASSOWRD,
	})
}

func (upy UPYun) Upload(storage string, localPath string) (bool, interface{}) {
	up := NewUpYun()
	err := up.Put(&upyun.PutObjectConfig{
		Path:      storage,
		LocalPath: localPath,
	})
	fileInfo, err := up.GetInfo(storage)
	if err != nil {
		return false, err
	}
	var bucket model.Bucket

	if dbError := database.POSTGRES.Where("bucket_type = ?", "upyun").First(&bucket).Error; dbError != nil {
		return false, err
	}

	var fileMessage model.FilesMessage
	fileMessage = model.FilesMessage{
		FilesMessageName: fileInfo.Name,
		FilesMessageKey:  fileInfo.ETag,
		FilesMessageSize: int(fileInfo.Size),
		FilesMessageURL:  fmt.Sprintf(config.UPYUNLink + fileInfo.Name),
		BucketID:         bucket.ID,
	}
	database.POSTGRES.Save(&fileMessage)
	var result = make(map[string]interface{})
	result["buckets"] = bucket
	result["fileMessagesBasicSerialize"] = fileMessage.BasicSerialize()
	result["fileInfo"] = fileInfo
	result["fileMessages"] = fileMessage
	return true, result
}

func (upy UPYun) GetOne(id uint, path string) (bool, interface{}) {
	up := NewUpYun()
	var fileMessage model.FilesMessage
	if dbError := database.POSTGRES.Where("id = ?", id).Or("files_message_name = ?", path).First(&fileMessage).Error; dbError != nil {
		return false, nil
	}

	var bucket model.Bucket
	if dbError := database.POSTGRES.Where("id = ?", fileMessage.BucketID).First(&bucket).Error; dbError != nil {
		return false, nil
	}
	fileInfo, _ := up.GetInfo(fileMessage.FilesMessageName)

	var result = make(map[string]interface{})
	result["buckets"] = bucket
	result["fileMessagesBasicSerialize"] = fileMessage.BasicSerialize()
	result["fileInfo"] = fileInfo
	result["fileMessages"] = fileMessage

	return true, result
}

type PathInfo struct {
	Name string
	Date time.Time
}

type PathInfos []PathInfo

func (upy UPYun) ListPath() (bool, interface{}) {
	up := NewUpYun()
	objsChan := make(chan *upyun.FileInfo, 10)
	go func() {
		err := up.List(&upyun.GetObjectsConfig{
			Path:        "/",
			ObjectsChan: objsChan,
		})
		log.Println(err)
	}()

	var paths PathInfos
	for obj := range objsChan {
		fmt.Println(obj.Name)
		var path PathInfo
		path.Name = obj.Name
		path.Date = obj.Time
		paths = append(paths, path)
	}
	return true, paths
}

func (upy UPYun) Usage() (string, error) {
	up := NewUpYun()
	usage, err := up.Usage()
	if err != nil {
		return "-1", err
	}
	return strconv.FormatFloat(float64(usage)/float64(1024*1024), 'f', 2, 32) + "MB", nil
}

func (upy UPYun) DeleteImages(id uint, path string) (bool, interface{}) {

	var fileMessage model.FilesMessage
	if dbError := database.POSTGRES.Where("id = ?", id).Or("files_message_name = ?", path).First(&fileMessage).Error; dbError != nil {
		return false, nil
	}

	up := NewUpYun()
	err := up.Delete(&upyun.DeleteObjectConfig{
		Path:  fileMessage.FilesMessageName,
		Async: true,
	})
	if err != nil {
		log.Println("upyun delete", err)
		return false, nil
	}
	database.POSTGRES.Delete(&fileMessage)
	return true, fileMessage.BasicSerialize()

}
