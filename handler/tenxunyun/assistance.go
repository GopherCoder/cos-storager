package tenxunyun

import (
	"context"
	"cos-storager/config"
	"cos-storager/model"
	"cos-storager/pkg/database"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/tencentyun/cos-go-sdk-v5"
)

type TenXun struct {
}

func (t TenXun) Upload(filePath string) (bool, interface{}) {
	ak := config.TXAK
	sk := config.TXSK
	bucket := config.TXBUCKET
	region := config.TXREGION
	bucketURL := fmt.Sprintf("http://%s.cos.%s.myqcloud.com", bucket, region)
	u, _ := url.Parse(bucketURL)
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Timeout: 100 * time.Second,
		Transport: &cos.AuthorizationTransport{
			SecretID:  ak,
			SecretKey: sk,
		},
	})
	f, err := os.Open(filePath)
	if err != nil {
		return false, nil
	}

	var oneBucket model.Bucket
	if dbError := database.POSTGRES.Where("bucket_url = ?", bucketURL).First(&oneBucket).Error; dbError != nil {
		return false, nil
	}

	opt := &cos.MultiUploadOptions{
		OptIni:   nil,
		PartSize: 1,
	}

	v, _, err := c.Object.MultiUpload(context.Background(), filePath, f, opt)
	if err != nil {
		fmt.Println(err)
		return false, nil
	}

	if err := os.Remove(filePath); err != nil {
		return false, nil
	}

	stat, _ := f.Stat()
	var oneMessage model.FilesMessage
	oneMessage = model.FilesMessage{
		FilesMessageName: v.Bucket,
		FilesMessageKey:  v.Key,
		FilesMessageURL:  fmt.Sprintf("http://%s", v.Location),
		FilesMessageSize: int(stat.Size()),
		BucketID:         oneBucket.ID,
	}
	database.POSTGRES.Save(&oneMessage)
	var result = make(map[string]interface{})
	result["buckets"] = oneBucket
	result["fileMessagesBasicSerialize"] = oneMessage.BasicSerialize()
	result["fileInfo"] = v
	result["fileMessages"] = oneMessage
	return true, result
}

func (t TenXun) GetOne(id uint) (bool, interface{}) {

	var oneFileMessage model.FilesMessage
	if dbError := database.POSTGRES.Where("id = ?", id).First(&oneFileMessage).Error; dbError != nil {
		return false, nil
	}

	var oneBucket model.Bucket
	if dbError := database.POSTGRES.Where("id = ?", oneFileMessage.BucketID).First(&oneBucket).Error; dbError != nil {
		return false, nil
	}

	var results = make(map[string]interface{})
	results["buckets"] = oneBucket.BasicSerializer()
	results["fileMessages"] = oneFileMessage.BasicSerialize()

	return true, results
}
func (t TenXun) GetAll() (bool, interface{}) {

	var oneBucket model.Bucket
	if dbError := database.POSTGRES.Where("bucket_type = ?", config.TXBUCKET).Or("bucket_type = ?", "tenxunyun").First(&oneBucket).Error; dbError != nil {
		return false, nil
	}

	var oneFileMessages model.FilesMessages
	if dbError := database.POSTGRES.Where("bucket_id = ?", oneBucket.ID).Find(&oneFileMessages).Error; dbError != nil {
		return false, nil
	}

	var results = make(map[string]interface{})
	results["buckets"] = oneBucket.BasicSerializer()
	var (
		messages model.Messages
	)
	for _, one := range oneFileMessages {
		messages = append(messages, one.BasicSerialize())
	}
	results["fileMessages"] = messages

	return true, results
}

func (t TenXun) DeleteOne(id uint) (bool, interface{}) {

	var oneFileMessage model.FilesMessage
	if dbError := database.POSTGRES.Where("id = ?", id).First(&oneFileMessage).Error; dbError != nil {
		return false, nil
	}

	var oneBucket model.Bucket
	if dbError := database.POSTGRES.Where("id = ?", oneFileMessage.BucketID).First(&oneBucket).Error; dbError != nil {
		return false, nil
	}

	ak := config.TXAK
	sk := config.TXSK
	bucket := config.TXBUCKET
	region := config.TXREGION
	bucketURL := fmt.Sprintf("http://%s.cos.%s.myqcloud.com", bucket, region)
	u, _ := url.Parse(bucketURL)
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Timeout: 100 * time.Second,
		Transport: &cos.AuthorizationTransport{
			SecretID:  ak,
			SecretKey: sk,
		},
	})

	_, err := c.Object.Delete(context.Background(), oneFileMessage.FilesMessageKey)
	if err != nil {
		return false, nil
	}

	database.POSTGRES.Delete(&oneFileMessage)

	var results = make(map[string]interface{})
	results["buckets"] = oneBucket.BasicSerializer()
	results["fileMessages"] = oneFileMessage.BasicSerialize()
	return true, results
}

func (t TenXun) DeleteAll() (bool, interface{}) {

	var oneBucket model.Bucket
	if dbError := database.POSTGRES.Where("bucket_type = ?", config.TXBUCKET).First(&oneBucket).Error; dbError != nil {
		return false, nil
	}

	var oneFileMessages model.FilesMessages
	if dbError := database.POSTGRES.Where("bucket_id = ?", oneBucket.ID).Find(&oneFileMessages).Error; dbError != nil {
		return false, nil
	}

	ak := config.TXAK
	sk := config.TXSK
	bucket := config.TXBUCKET
	region := config.TXREGION
	bucketURL := fmt.Sprintf("http://%s.cos.%s.myqcloud.com", bucket, region)
	u, _ := url.Parse(bucketURL)
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Timeout: 100 * time.Second,
		Transport: &cos.AuthorizationTransport{
			SecretID:  ak,
			SecretKey: sk,
		},
	})

	var names []string
	var fileMessageResult model.Messages
	for _, one := range oneFileMessages {
		names = append(names, one.FilesMessageKey)
		fileMessageResult = append(fileMessageResult, one.BasicSerialize())
	}

	ctx := context.Background()
	obs := []cos.Object{}
	for _, v := range names {
		obs = append(obs, cos.Object{Key: v})
	}
	opt := &cos.ObjectDeleteMultiOptions{
		Objects: obs,
	}

	v, _, err := c.Object.DeleteMulti(ctx, opt)
	if err != nil {
		return false, nil
	}

	var results = make(map[string]interface{})
	results["buckets"] = oneBucket.BasicSerializer()
	results["fileMessages"] = fileMessageResult
	results["fileInfo"] = v

	return true, results
}
