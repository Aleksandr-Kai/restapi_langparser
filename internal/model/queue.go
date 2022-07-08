package model

import (
	"gorm.io/gorm"
	"time"
)

type Queue struct {
	gorm.Model
	DomainID uint `gorm:"unique"`
	UpdateAt time.Time
}
