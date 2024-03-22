package types

import "gorm.io/gorm"

type DictType string

const (
	DictTypeString DictType = "string"
	DictTypeInt    DictType = "int"
	DictTypeBool   DictType = "bool"
)

func (t DictType) String() string {
	return string(t)
}

// SysDict 字典配置表
type SysDict struct {
	gorm.Model
	ParentID       uint      `gorm:"column:parent_id;type:bigint(20);NOT NULL"`        // 父ID
	Code           string    `gorm:"column:code;type:varchar(32);NOT NULL"`            // 字典编号
	DictValue      string    `gorm:"column:dict_value;type:varchar(128);NOT NULL"`     // 字典值（接口传参）
	DictLabel      string    `gorm:"column:dict_label;type:varchar(128);NOT NULL"`     // 字典名称
	DictType       string    `gorm:"column:dict_type;type:varchar(10);default:string"` // 字典值类型（string、int、bool）
	Sort           int       `gorm:"column:sort;type:int(11)"`                         // 排序
	Remark         string    `gorm:"column:remark;type:varchar(255)"`                  // 字典备注
	ParentDictType string    `gorm:"-"`                                                // 父字典类型
	Children       []SysDict `gorm:"-"`                                                // 不映射到数据库，用于存储子项
}

func (m *SysDict) TableName() string {
	return "sys_dict"
}
