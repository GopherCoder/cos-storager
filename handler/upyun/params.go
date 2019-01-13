package upyun

import "mime/multipart"

type GetUpYunOneParams struct {
	ID   uint   `form:"id" json:"id"`
	Path string `form:"path" json:"path"`
}

type GetUpYunByPathParams struct {
	Path string `form:"path" json:"path"`
	//Limit  string `form:"limit,default:10" json:"limit"`
	//Offset string `form:"offset,default:10" json:"offset"`
}

type PostUpYunParams struct {
	Path string         `form:"path" json:"path"`
	File multipart.Form `form:"file_path" json:"file" binding:"required"`
}

type DeleteUpYunParams struct {
	ID   uint   `form:"id" json:"id"`
	Path string `form:"path" json:"path"`
}
