package types

import "gorm.io/gorm"

// SysAudit 审计日志
type SysAudit struct {
	gorm.Model
	Operator      string  `gorm:"column:operator;type:varchar(100);NOT NULL"`         // 操作人
	TenantID      uint    `gorm:"column:tenant_id;type:bigint(20) unsigned;NOT NULL"` // 租户ID
	RequestMethod string  `gorm:"column:request_method;type:varchar(20);NOT NULL"`    // 请求方法
	RequestUrl    string  `gorm:"column:request_url;type:varchar(255);NOT NULL"`      // 请求地址
	RequestBody   string  `gorm:"column:request_body;type:json"`                      // 请求body体
	ResponseBody  string  `gorm:"column:response_body;type:json"`                     // 响应内容
	IsError       bool    `gorm:"column:is_error;type:tinyint(1);NOT NULL"`           // 错误
	ErrorMessage  string  `gorm:"column:error_message;type:text"`                     // 错误信息
	TraceID       string  `gorm:"column:trace_id;type:varchar(100);NOT NULL"`         // 追踪ID
	Duration      float64 `gorm:"column:duration;type:double;NOT NULL"`               // 响应时间(单位秒)
}

func (m *SysAudit) TableName() string {
	return "sys_audit"
}
