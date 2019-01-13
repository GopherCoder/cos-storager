package qiniuyun

import "mime/multipart"

type PostUploadParams struct {
	FilePath multipart.Form `json:"file_path"`
}
