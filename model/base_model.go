package model

func GetAllModels() []interface{} {
	return []interface{}{
		&Bucket{},
		&FilesMessage{},
	}
}
