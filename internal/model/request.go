package model

import (
	"gorm.io/gorm"
	"time"
)

type Request struct {
	DomainID  uint   `gorm:"primaryKey"`
	Code      string `gorm:"primaryKey"`
	CreatedAt time.Time
}

func (r *Request) BeforeCreate(*gorm.DB) (err error) {
	r.CreatedAt = time.Now()
	return nil
}
