package model

import "gorm.io/gorm"

type Queue struct {
	gorm.Model
	DomainID uint `gorm:"unique"` //`gorm:"primaryKey"`
}
