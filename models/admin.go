package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (admin *Admin) BeforeCreate(tx *gorm.DB) (err error) {
	admin.ID = uuid.New()
	return
}

type Admin struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`
	Name      string    `json:"name" gorm:"type:varchar(255);not null"`
	Email     string    `json:"email" gorm:"type:varchar(255);uniqueIndex;not null"`
	Password  string    `json:"password" gorm:"type:varchar(255);not null"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
	Products  []Product `json:"products" gorm:"foreignKey:AdminID"`
}
