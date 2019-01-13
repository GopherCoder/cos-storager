package model

import "github.com/jinzhu/gorm"

type Bucket struct {
	gorm.Model
	BucketName string `gorm:"type:varchar" json:"bucket_name"`
	BucketType string `gorm:"type:varchar" json:"bucket_type"`
	BucketURL  string `gorm:"type:varchar" json:"bucket_url"`
}

func (Bucket) TableName() string {
	return "buckets"
}

type Buckets []Bucket

type BucketMessage struct {
	Name string `json:"name"`
	Type string `json:"type"`
	URL  string `json:"url"`
	ID   uint   `json:"id"`
}

func (b Bucket) BasicSerializer() BucketMessage {
	return BucketMessage{
		Name: b.BucketName,
		Type: b.BucketType,
		URL:  b.BucketURL,
		ID:   b.ID,
	}
}
