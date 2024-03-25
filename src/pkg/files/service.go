package files

import (
	"context"
	"crypto/md5"
	"embed"
	"encoding/hex"
	"fmt"
	"github.com/IceBearAI/aigc/src/repository"
	"github.com/IceBearAI/aigc/src/repository/types"
	"github.com/IceBearAI/aigc/src/services"
	"github.com/IceBearAI/aigc/src/util"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"io"
	"mime/multipart"
	"os"
	"path"
	"strings"
	"time"
)

// CreationOptions is the options for the faceswap service.
type CreationOptions struct {
	httpClientOpts []kithttp.ClientOption
	storageType    string
	serverUrl      string
	localDataPath  string
	s3Endpoint     string
	s3AccessKey    string
	s3SecretKey    string
	s3BucketName   string
	s3Region       string
}

// CreationOption is a creation option for the faceswap service.
type CreationOption func(*CreationOptions)

// WithHTTPClientOptions returns a CreationOption that sets the http client options.
func WithHTTPClientOptions(opts ...kithttp.ClientOption) CreationOption {
	return func(o *CreationOptions) {
		o.httpClientOpts = opts
	}
}

// WithStorageType returns a CreationOption that sets the storage type.
func WithStorageType(storageType string) CreationOption {
	return func(o *CreationOptions) {
		o.storageType = storageType
	}
}

// WithServerUrl returns a CreationOption that sets the server url.
func WithServerUrl(serverUrl string) CreationOption {
	return func(o *CreationOptions) {
		o.serverUrl = serverUrl
	}
}

// WithLocalDataPath returns a CreationOption that sets the local data path.
func WithLocalDataPath(localDataPath string) CreationOption {
	return func(o *CreationOptions) {
		o.localDataPath = localDataPath
	}
}

// WithS3Endpoint returns a CreationOption that sets the s3 endpoint.
func WithS3Endpoint(s3Endpoint string) CreationOption {
	return func(o *CreationOptions) {
		o.s3Endpoint = s3Endpoint
	}
}

// WithS3AccessKey returns a CreationOption that sets the s3 access key.
func WithS3AccessKey(s3AccessKey string) CreationOption {
	return func(o *CreationOptions) {
		o.s3AccessKey = s3AccessKey
	}
}

// WithS3SecretKey returns a CreationOption that sets the s3 secret key.
func WithS3SecretKey(s3SecretKey string) CreationOption {
	return func(o *CreationOptions) {
		o.s3SecretKey = s3SecretKey
	}
}

// WithS3BucketName returns a CreationOption that sets the s3 bucket name.
func WithS3BucketName(s3BucketName string) CreationOption {
	return func(o *CreationOptions) {
		o.s3BucketName = s3BucketName
	}
}

type Service interface {
	CreateFile(ctx context.Context, request FileRequest) (file File, err error)
	ListFiles(ctx context.Context, request ListFileRequest) (files FileList, err error)
	GetFile(ctx context.Context, fileId string) (file File, err error)
	DeleteFile(ctx context.Context, fileId string) (err error)
	// UploadToS3 上传文件到s3存储
	UploadToS3(ctx context.Context, file multipart.File, fileType string, isPublicBucket bool) (s3Url string, err error)
	// UploadLocal 上传文件到本地存储
	UploadLocal(ctx context.Context, file multipart.File, fileType string) (localFile string, err error)
	// UploadToStorage 上传文件到存储
	UploadToStorage(ctx context.Context, file multipart.File, fileType string) (url string, err error)
}

type service struct {
	logger      log.Logger
	traceId     string
	store       repository.Repository
	apiSvc      services.Service
	localDataFS embed.FS
	options     *CreationOptions
}

func (s *service) UploadToStorage(ctx context.Context, file multipart.File, fileType string) (url string, err error) {
	switch s.options.storageType {
	case "s3":
		url, err = s.UploadToS3(ctx, file, fileType, false)
		break
	case "local":
		url, err = s.UploadLocal(ctx, file, fileType)
		break
	}
	return
}

func (s *service) UploadToS3(ctx context.Context, file multipart.File, fileType string, isPublicBucket bool) (s3Url string, err error) {
	// 如果 isPublicBucket 为true,则使用公有桶，否则使用私有桶
	//bucketName := s.options.s3BucketName //默认私有桶
	//if isPublicBucket == true {
	//	//bucketName = s.options.s3BucketNamePublic
	//}
	//logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "UploadToS3")
	//paths := strings.Split(time.Now().Format(time.RFC3339), "-")
	//id, _ := util.GenShortId(24)
	//fileName := fmt.Sprintf("%s.%s", id, fileType)
	//targetPath := fmt.Sprintf("%s/%s/%s/%s", fileType, paths[0], paths[1], fileName)
	//err = s.apiSvc.S3Client(ctx).Upload(ctx, bucketName, targetPath, file, "")
	//if err != nil {
	//	return "", err
	//}
	//
	//s3Url, err = s.apiSvc.S3Client(ctx).ShareGen(ctx, bucketName, targetPath, 60*24*31*12)
	//if err != nil {
	//	_ = level.Error(logger).Log("pan", "ShareGen", "err", err.Error())
	//	return
	//}

	return s3Url, nil
}

// UploadLocal 将文件上传到本地目录
func (s *service) UploadLocal(ctx context.Context, file multipart.File, fileType string) (localFile string, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "service", "UploadLocal")
	paths := strings.Split(time.Now().Format(time.RFC3339), "-")
	id, _ := util.GenShortId(24)
	fileName := fmt.Sprintf("%s.%s", id, fileType)
	targetPath := path.Join(s.options.localDataPath, fileType, paths[0], paths[1], fileName)
	if err = os.MkdirAll(path.Join(s.options.localDataPath, fileType, paths[0], paths[1]), os.ModePerm); err != nil {
		return "", errors.Wrap(err, "os.MkdirAll")
	}

	// 创建一个本地文件，用于保存上传的文件
	dst, err := os.Create(targetPath)
	if err != nil {
		_ = level.Error(logger).Log("os", "Create", "err", err.Error())
		return
	}
	defer func(dst *os.File) {
		_ = dst.Close()
	}(dst)

	// 将上传的文件复制到本地文件
	if _, err = io.Copy(dst, file); err != nil {
		_ = level.Error(logger).Log("io", "Copy", "err", err.Error())
		return
	}
	return fmt.Sprintf("%s/%s", s.options.serverUrl, path.Join(fileType, paths[0], paths[1], fileName)), nil
}

func (s *service) CreateFile(ctx context.Context, request FileRequest) (file File, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "CreateFile")

	defer request.File.Close()
	// 计算文件md5
	hash := md5.New()
	if _, err = io.Copy(hash, request.File); err != nil {
		_ = level.Error(logger).Log("io.Copy", "err", err.Error())
		return file, err
	}
	md5Str := hex.EncodeToString(hash.Sum(nil))
	if _, err = request.File.Seek(0, io.SeekStart); err != nil {
		_ = level.Error(logger).Log("request.file.Seek", "err", err.Error())
		return file, err
	}

	// 根据md5查询文件是否已经存在
	if res, err := s.store.Files().FindFileByMd5(ctx, md5Str); err == nil && res.ID > 0 {
		file = File{
			FileId:    res.FileID,
			FileName:  res.Name,
			Size:      res.Size,
			FileType:  res.Type,
			S3Url:     res.S3Url,
			Purpose:   res.Purpose,
			TenantId:  res.TenantID,
			CreatedAt: res.CreatedAt,
			LineCount: res.LineCount,
		}
		return file, nil
	}

	fileUrl, err := s.UploadToStorage(ctx, request.File, request.FileType)
	if err != nil {
		_ = level.Error(logger).Log("uploadLocal", err.Error())
		return
	}
	_ = level.Info(logger).Log("fileUrl", fileUrl)
	// 保存文件信息到数据库
	data := &types.Files{
		FileID:     uuid.New().String(),
		Name:       request.Header.Filename,
		Size:       request.Header.Size,
		Type:       request.FileType,
		Md5:        md5Str,
		S3Url:      fileUrl,
		Purpose:    request.Purpose,
		TenantID:   request.TenantId,
		LineCount:  request.LineCount,
		TokenCount: request.TokenCount,
	}
	err = s.store.Files().CreateFile(ctx, data)
	if err != nil {
		_ = level.Error(logger).Log("store.Chat", "CreateFile", "err", err.Error())
		return file, err
	}
	file = File{
		FileId:    data.FileID,
		FileName:  data.Name,
		Size:      data.Size,
		FileType:  data.Type,
		S3Url:     data.S3Url,
		Purpose:   data.Purpose,
		TenantId:  data.TenantID,
		CreatedAt: data.CreatedAt,
		LineCount: data.LineCount,
	}
	return file, nil
}

func (s *service) ListFiles(ctx context.Context, request ListFileRequest) (files FileList, err error) {
	files.Files = make([]File, 0)
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "ListFiles")
	res, total, err := s.store.Files().ListFiles(ctx, request.TenantId, request.Purpose, request.FileName, request.FileType, request.Page, request.PageSize)
	if err != nil {
		_ = level.Error(logger).Log("store.Chat", "ListFiles", "err", err.Error())
		return files, err
	}
	for _, v := range res {
		files.Files = append(files.Files, File{
			FileId:    v.FileID,
			FileName:  v.Name,
			Size:      v.Size,
			FileType:  v.Type,
			S3Url:     v.S3Url,
			Purpose:   v.Purpose,
			TenantId:  v.TenantID,
			CreatedAt: v.CreatedAt,
		})
	}
	files.Total = total
	return files, nil
}

func (s *service) GetFile(ctx context.Context, fileId string) (file File, err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "GetFile")
	res, err := s.store.Files().FindFileByFileId(ctx, fileId)
	if err != nil {
		_ = level.Error(logger).Log("store.Chat", "FindFileByFileId", "err", err.Error(), "fileId", fileId)
		return file, err
	}
	file = File{
		FileId:    res.FileID,
		FileName:  res.Name,
		Size:      res.Size,
		FileType:  res.Type,
		S3Url:     res.S3Url,
		Purpose:   res.Purpose,
		TenantId:  res.TenantID,
		CreatedAt: res.CreatedAt,
	}
	return file, nil
}

func (s *service) DeleteFile(ctx context.Context, fileId string) (err error) {
	logger := log.With(s.logger, s.traceId, ctx.Value(s.traceId), "method", "DeleteFile")
	err = s.store.Files().DeleteFile(ctx, fileId)
	if err != nil {
		_ = level.Error(logger).Log("store.Chat", "DeleteFile", "err", err.Error(), "fileId", fileId)
		return err
	}
	return nil
}

type s3Config struct {
	AccessKey        string
	SecretKey        string
	BucketName       string //默认私有桶配置
	BucketNamePublic string //公有桶配置
	ProjectName      string
}

type Config struct {
	S3            s3Config
	LocalDataPath string
	ServerUrl     string
}

func NewService(logger log.Logger, traceId string, store repository.Repository, apiSvc services.Service, opts ...CreationOption) Service {
	_ = log.With(logger, "pkg.files", "service")
	options := &CreationOptions{
		storageType:   "local",
		serverUrl:     "http://localhost:8080",
		localDataPath: "/data/storage",
	}
	for _, opt := range opts {
		opt(options)
	}
	return &service{
		logger:  logger,
		traceId: traceId,
		store:   store,
		apiSvc:  apiSvc,
		options: options,
	}
}
