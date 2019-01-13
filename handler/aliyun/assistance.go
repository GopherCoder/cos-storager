package aliyun

import (
	"cos-storager/config"
	"cos-storager/model"
	"venus/pkg/database"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type AliYun struct {
	EndPoint        string
	Bucket          string
	AccessKeyID     string
	AccessKeySecret string
}

func NewAliYun() *AliYun {
	return &AliYun{
		EndPoint:        config.ALIYUNENDPOINT,
		Bucket:          config.ALIYUNBUCKET,
		AccessKeyID:     config.ALIYUNACCESSKEYID,
		AccessKeySecret: config.ALIYUNACCESSKEYSECRET,
	}
}

func Client() *oss.Client {
	aliYun := NewAliYun()
	aliYunClient, err := oss.New(aliYun.EndPoint, aliYun.AccessKeyID, aliYun.AccessKeySecret)
	if err != nil {
		return nil
	}
	return aliYunClient
}

func (aa AliYun) Upload(fileName string, localPath string) (bool, interface{}) {
	var bucketType model.Bucket
	if dbError := database.POSTGRES.Where("bucket_type = ?", config.ALIYUNTYPE).First(&bucketType).Error; dbError != nil {
		return false, nil
	}
	client := Client()
	bucket, err := client.Bucket(aa.Bucket)
	if err != nil {
		return false, nil
	}

	err = bucket.PutObjectFromFile(fileName, localPath)
	if err != nil {
		return false, nil
	}
	//bucket.GetObject()
	return true, nil

}
