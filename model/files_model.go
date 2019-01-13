package model

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type FilesMessage struct {
	gorm.Model
	FilesMessageName string `gorm:"type:varchar" json:"files_message_name"`
	FilesMessageURL  string `gorm:"type:varchar" json:"files_message_url"`
	FilesMessageKey  string `gorm:"type:varchar" json:"files_message_key"`
	FilesMessageSize int    `gorm:"type:integer" json:"files_message_size"`
	BucketID         uint
}

func (FilesMessage) TableName() string {
	return "files_messages"
}

type FilesMessages []FilesMessage

type Message struct {
	OriginURL   string `json:"origin_url"`
	HtmlURL     string `json:"html_url"`
	UbbURL      string `json:"ubb_url"`
	MarkDownURL string `json:"mark_down_url"`
}

type Messages []Message

func (f FilesMessage) BasicSerialize() Message {
	return Message{
		OriginURL:   f.FilesMessageURL,
		HtmlURL:     fmt.Sprintf(`<img src="%s"/>`, f.FilesMessageURL),
		UbbURL:      fmt.Sprintf(`[IMG]%s[/IMG]`, f.FilesMessageURL),
		MarkDownURL: fmt.Sprintf(`![](%s)`, f.FilesMessageURL),
	}
}
