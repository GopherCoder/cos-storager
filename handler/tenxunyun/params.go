package tenxunyun

import "mime/multipart"

type PostParams struct {
	FilePath multipart.Form `form:"file_path" json:"file_path"`
}
