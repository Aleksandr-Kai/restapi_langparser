package model

type TagsLangs struct {
	ID       int64  `gorm:"primaryKey;column:id" json:"-"`
	DomainID int64  `gorm:"column:domain_id" json:"-"`
	Lang     string `gorm:"column:lang"`
}
