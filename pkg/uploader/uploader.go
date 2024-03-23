package uploader

import (
	"mime/multipart"
)

type ImageUploader interface {
	Upload(folder string, file *multipart.FileHeader) (string, error)
}
