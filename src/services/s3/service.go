package s3

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pkg/errors"
	"mime/multipart"
	"time"
)

type Middleware func(Service) Service

type Service interface {
	Upload(ctx context.Context, bucketName string, targetPath string, file multipart.File, header string) (err error)
	// ShareGen 生成分享链接
	// expireMinute: 分钟
	ShareGen(ctx context.Context, bucketName, targetPath string, expireMinute int64) (url string, err error)
}

type service struct {
	client      *s3.S3
	downloadUrl string
	session     *session.Session
}

func (s *service) ShareGen(ctx context.Context, bucketName, targetPath string, expireMinute int64) (url string, err error) {
	req, _ := s.client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(targetPath),
	})
	url, err = req.Presign(time.Duration(expireMinute) * time.Minute)
	if err != nil {
		return "", errors.Wrap(err, "s3.ShareGen")
	}
	return
}

func (s *service) Upload(
	ctx context.Context,
	bucketName, targetPath string,
	file multipart.File,
	header string,
) (err error) {
	uploader := s3manager.NewUploader(s.session)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(targetPath),
		Body:        file,
		ContentType: aws.String(header),
	}, func(u *s3manager.Uploader) {
		u.PartSize = 10 * 1024 * 1024
		u.LeavePartsOnError = true
		u.Concurrency = 3
	})
	if err != nil {
		return errors.Wrap(err, "s3.Upload")
	}
	return
}

func New(storageType, accessKey, secretKey, s3url, region string) Service {
	if storageType == "local" {
		return NewLocal(s3url)
	}
	s3Config := &aws.Config{
		Credentials:                    credentials.NewStaticCredentials(accessKey, secretKey, ""),
		Endpoint:                       aws.String(s3url),
		Region:                         aws.String(region),
		DisableSSL:                     aws.Bool(true),
		S3ForcePathStyle:               aws.Bool(true), //virtual-host style方式，不要修改
		DisableRestProtocolURICleaning: aws.Bool(true),
	}
	newSession, _ := session.NewSession(s3Config)
	s3Client := s3.New(newSession)
	return &service{client: s3Client, session: newSession}
}
