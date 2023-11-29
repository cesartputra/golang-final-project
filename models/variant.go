package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (variant *Variant) BeforeCreate(tx *gorm.DB) (err error) {
	variant.ID = uuid.New()
	return
}

type Variant struct {
	ID          uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`
	VariantName string    `json:"variantName" gorm:"type:varchar(255);not null"`
	Quantity    int       `json:"quantity" gorm:"type:integer;not null"`
	ProductID   uuid.UUID `json:"productID" gorm:"type:char(36);not null"`
	CreatedAt   time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}
