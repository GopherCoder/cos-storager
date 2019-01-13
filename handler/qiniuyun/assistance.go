package qiniuyun

import (
	"context"
	"cos-storager/config"
	"cos-storager/model"
	"cos-storager/pkg/database"
	"fmt"
	"strings"

	"qiniupkg.com/x/log.v7"

	"github.com/qiniu/api.v7/auth/qbox"

	"github.com/qiniu/api.v7/storage"
)

type QiNiu struct {
}

type QiNiuResult struct {
	Key    string `json:"key"`
	Hash   string `json:"hash"`
	Fsize  int    `json:"Fsize"`
	Bucket string `json:"bucket"`
	Name   string `json:"name"`
}

func (q QiNiu) Upload(filePath string) (bool, interface{}) {
	bucket := config.QNBUCKET
	accessKey := config.QNAK
	secretKey := config.QNSK
	tempList := strings.Split(filePath, "/")
	if len(tempList) == 0 {
		return false, "-1"
	}
	key := tempList[len(tempList)-1]
	putPolicy := storage.PutPolicy{
		Scope:      bucket,
		ReturnBody: `{"key":"$(key)","hash":"$(etag)","fsize":$(fsize),"bucket":"$(bucket)","name":"$(x:name)"}`,
	}
	mac := qbox.NewMac(accessKey, secretKey)
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{}
	formUploader := storage.NewFormUploader(&cfg)
	ret := QiNiuResult{}
	putExtra := storage.PutExtra{
		Params: map[string]string{
			"x:name": "cos-storage",
		},
	}
	err := formUploader.PutFile(context.Background(), &ret, upToken, key, filePath, &putExtra)
	if err != nil {
		log.Println(err)
		return false, "-1"
	}
	fmt.Println(ret)
	return true, ret
}

func (q QiNiu) GetOne(id uint) (bool, interface{}) {
	accessKey := config.QNAK
	secretKey := config.QNSK
	mac := qbox.NewMac(accessKey, secretKey)
	cfg := storage.Config{}
	bucketManager := storage.NewBucketManager(mac, &cfg)

	var (
		oneMessage model.FilesMessage
		oneBucket  model.Bucket
	)
	if notFound := database.POSTGRES.Where("id = ?", id).First(&oneMessage).RecordNotFound(); notFound {
		return false, "-1"
	}
	if notFound := database.POSTGRES.Where("id = ?", oneMessage.BucketID).First(&oneBucket).RecordNotFound(); notFound {
		return false, "-1"
	}
	fileInfo, err := bucketManager.Stat(oneBucket.BucketType, oneMessage.FilesMessageKey)
	if err != nil {
		return false, "-1"
	}

	var result = make(map[string]interface{})
	result["bucket"] = oneBucket
	result["fileMessage"] = oneMessage
	result["fileInfo"] = fileInfo
	return true, result
}

func (q QiNiu) GetAll() (bool, interface{}) {
	var (
		oneMessages model.FilesMessages
		oneBuckets  model.Buckets
	)
	if dbError := database.POSTGRES.Where("bucket_type = ?", "qiniuyun").Find(&oneBuckets).Error; dbError != nil {
		return false, nil
	}
	var oneBucketsIDs []uint
	for _, bucket := range oneBuckets {
		oneBucketsIDs = append(oneBucketsIDs, bucket.ID)
	}
	if dbError := database.POSTGRES.Where("bucket_id in (?)", oneBucketsIDs).Find(&oneMessages).Error; dbError != nil {
		return false, nil
	}

	accessKey := config.QNAK
	secretKey := config.QNSK
	mac := qbox.NewMac(accessKey, secretKey)
	cfg := storage.Config{}
	bucketManager := storage.NewBucketManager(mac, &cfg)

	var keys []string
	for _, key := range oneMessages {
		keys = append(keys, key.FilesMessageKey)
	}
	statOps := make([]string, 0, len(keys))
	for _, key := range keys {
		statOps = append(statOps, storage.URIStat(oneBuckets[0].BucketType, key))
	}
	results, err := bucketManager.Batch(statOps)
	if err != nil {
		return false, nil
	}

	var allResults = make(map[string]interface{})
	allResults["fileMessages"] = oneMessages
	allResults["buckets"] = oneBuckets
	allResults["fileInfo"] = results

	return true, allResults
}

func (q QiNiu) DeleteOne(id uint) (bool, interface{}) {
	accessKey := config.QNAK
	secretKey := config.QNSK
	mac := qbox.NewMac(accessKey, secretKey)
	cfg := storage.Config{}
	bucketManager := storage.NewBucketManager(mac, &cfg)

	var oneMessages model.FilesMessage
	if dbError := database.POSTGRES.Where("id = ?", id).First(&oneMessages).Error; dbError != nil {
		return false, nil
	}

	var oneBucket model.Bucket
	if dbError := database.POSTGRES.Where("id = ?", oneMessages.BucketID).First(&oneBucket).Error; dbError != nil {
		return false, nil
	}

	err := bucketManager.Delete(oneBucket.BucketType, oneMessages.FilesMessageKey)
	if err != nil {
		return false, nil
	}

	if oneMessages.ID != 0 {
		database.POSTGRES.Delete(&oneMessages)
	}

	var results = make(map[string]interface{})
	results["buckets"] = oneBucket
	results["fileMessages"] = oneMessages

	return true, results
}

func (q QiNiu) DeleteAll() (bool, interface{}) {
	var (
		oneMessages model.FilesMessages
		oneBucket   model.Bucket
	)

	if dbError := database.POSTGRES.Where("bucket_type = ?", "qiniuyun").First(&oneBucket).Error; dbError != nil {
		return false, nil
	}

	if dbError := database.POSTGRES.Where("bucket_id = ?", oneBucket.ID).Find(&oneMessages).Error; dbError != nil {
		return false, nil
	}

	deleteOps := make([]string, 0, len(oneMessages))
	for _, key := range oneMessages {
		deleteOps = append(deleteOps, storage.URIDelete(oneBucket.BucketType, key.FilesMessageKey))
	}
	accessKey := config.QNAK
	secretKey := config.QNSK
	mac := qbox.NewMac(accessKey, secretKey)
	cfg := storage.Config{}
	bucketManager := storage.NewBucketManager(mac, &cfg)
	results, err := bucketManager.Batch(deleteOps)
	if err != nil {
		return false, nil
	}

	for _, one := range oneMessages {
		if one.ID != 0 {
			database.POSTGRES.Delete(&one)
		}
	}
	var allResults = make(map[string]interface{})
	allResults["buckets"] = oneBucket
	allResults["fileMessages"] = oneMessages
	allResults["fileInfos"] = results

	return true, allResults
}
