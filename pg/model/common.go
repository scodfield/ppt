package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        string `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt int64  `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt int64  `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (b *BaseModel) BeforeCreate(db *gorm.DB) error {
	uuidV7, err := uuid.NewV7()
	if err != nil {
		return err
	}
	b.ID = uuidV7.String()
	return nil
}
