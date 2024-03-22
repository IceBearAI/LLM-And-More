package files

import (
	"context"
	"github.com/IceBearAI/aigc/src/helpers/page"
	"github.com/IceBearAI/aigc/src/repository/types"
	"gorm.io/gorm"
)

type Middleware func(service Service) Service

type Service interface {
	// CreateFile 创建文件
	CreateFile(ctx context.Context, data *types.Files) (err error)
	// FindFileByFileId 根据fileId查找文件
	FindFileByFileId(ctx context.Context, fileId string) (res types.Files, err error)
	// FindFileByMd5 根据md5查找文件
	FindFileByMd5(ctx context.Context, md5 string) (res types.Files, err error)
	// ListFiles 文件分页列表
	ListFiles(ctx context.Context, tenantId uint, purpose string, fileName string, fileType string, page, pageSize int) (res []types.Files, total int64, err error)
	// DeleteFile 删除文件
	DeleteFile(ctx context.Context, fileId string) (err error)
}

type service struct {
	db *gorm.DB
}

func (s *service) FindFileByMd5(ctx context.Context, md5 string) (res types.Files, err error) {
	err = s.db.WithContext(ctx).Where("md5 = ?", md5).First(&res).Error
	return
}

func (s *service) CreateFile(ctx context.Context, data *types.Files) (err error) {
	err = s.db.WithContext(ctx).Create(data).Error
	return
}

func (s *service) FindFileByFileId(ctx context.Context, fileId string) (res types.Files, err error) {
	err = s.db.WithContext(ctx).Where("file_id = ?", fileId).First(&res).Error
	return
}

func (s *service) ListFiles(ctx context.Context, tenantId uint, purpose string, fileName string, fileType string, pageNum, pageSize int) (res []types.Files, total int64, err error) {
	query := s.db.WithContext(ctx).Model(&types.Files{})
	if tenantId > 0 {
		query = query.Where("tenant_id = ?", tenantId)
	}
	if purpose != "" {
		query = query.Where("purpose = ?", purpose)
	}
	if fileName != "" {
		query = query.Where("name LIKE ?", "%"+fileName+"%")
	}
	if fileType != "" {
		query = query.Where("type = ?", fileType)
	}
	limit, offset := page.Limit(pageNum, pageSize)
	err = query.Count(&total).Order("id DESC").Limit(limit).Offset(offset).Find(&res).Error
	return
}

func (s *service) DeleteFile(ctx context.Context, fileId string) (err error) {
	return s.db.WithContext(ctx).Where("file_id = ?", fileId).Delete(&types.Files{}).Error
}

func New(db *gorm.DB) Service {
	return &service{db: db}
}
