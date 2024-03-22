package types

import "gorm.io/gorm"

type FilePurpose string

const (
	FilePurposeFineTune     FilePurpose = "fine-tune"
	FilePurposeFineTuneEval FilePurpose = "fine-tune-eval"
)

func (f FilePurpose) String() string {
	return string(f)
}

type Files struct {
	gorm.Model
	FileID     string `gorm:"column:file_id;type:varchar(64);NOT NULL"`    // 文件ID
	Name       string `gorm:"column:name;type:varchar(64);NOT NULL"`       // 文件名称
	Size       int64  `gorm:"column:size;type:bigint(20);NOT NULL"`        // 文件大小
	Type       string `gorm:"column:type;type:varchar(64);NOT NULL"`       // 文件类型
	Md5        string `gorm:"column:md5;type:varchar(64);NOT NULL"`        // 文件md5
	S3Url      string `gorm:"column:s3_url;type:varchar(500)"`             // s3链接
	Purpose    string `gorm:"column:purpose;type:varchar(32);NOT NULL"`    // 文件的预期目的, 默认 fine-tune
	TenantID   uint   `gorm:"column:tenant_id;type:bigint(20);NOT NULL"`   // 租户ID
	LineCount  int    `gorm:"column:line_count;type:bigint(20);NOT NULL"`  // 文件行数
	TokenCount int    `gorm:"column:token_count;type:bigint(20);NOT NULL"` // 文件token数
}

func (m *Files) TableName() string {
	return "files"
}
