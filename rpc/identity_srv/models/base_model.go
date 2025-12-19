package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uuid.UUID      `gorm:"column:id;primaryKey;type:uuid;default:gen_random_uuid();comment:主键"`
	CreatedAt int64          `gorm:"column:created_at;autoCreateTime:milli;index;comment:创建时间"`
	UpdatedAt int64          `gorm:"column:updated_at;autoUpdateTime:milli;index;comment:更新时间"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index;comment:删除时间"`
}
