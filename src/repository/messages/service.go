package messages

import (
	"context"
	"github.com/IceBearAI/aigc/src/repository/types"
	"gorm.io/gorm"
)

type Service interface {
	Create(ctx context.Context, message *types.ChatMessages) (err error)
	Update(ctx context.Context, message *types.ChatMessages) (err error)
}

type service struct {
	db *gorm.DB
}

func (s *service) Update(ctx context.Context, message *types.ChatMessages) (err error) {
	err = s.db.WithContext(ctx).Model(message).Updates(message).Error
	return
}

func (s *service) Create(ctx context.Context, message *types.ChatMessages) (err error) {
	err = s.db.WithContext(ctx).Create(message).Error
	return
}

func New(db *gorm.DB) Service {
	return &service{db: db}
}
