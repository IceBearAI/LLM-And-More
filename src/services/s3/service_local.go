package s3

import (
	"context"
	"mime/multipart"
)

type serviceLocal struct {
	localPath   string
	downloadUrl string
}

func (s *serviceLocal) ShareGen(ctx context.Context, bucketName, targetPath string, expireMinute int64) (url string, err error) {
	return
}

func (s *serviceLocal) Upload(ctx context.Context, bucketName string, targetPath string, file multipart.File, header string) (err error) {
	// 将文件保存到本地

	return
}

func NewLocal(localPath string) Service {
	return &serviceLocal{localPath: localPath}
}
