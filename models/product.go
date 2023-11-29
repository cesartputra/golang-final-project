package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (product *Product) BeforeCreate(tx *gorm.DB) (err error) {
	product.ID = uuid.New()
	return
}

type Product struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`
	Name      string    `json:"name" gorm:"type:varchar(255);not null"`
	ImageUrl  string    `json:"imageUrl" gorm:"type:varchar(255);not null"`
	AdminID   uuid.UUID `json:"adminID" gorm:"type:char(36);not null"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
	Variants  []Variant `json:"variants" gorm:"foreignKey:ProductID"`
}
