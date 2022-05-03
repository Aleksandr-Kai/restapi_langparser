package model

type Result struct {
	RequestCode string `json:"request_code" gorm:"primaryKey;column:request_code"`
	DomainID    int    `json:"domain_id" gorm:"primaryKey;column:domain_id"`
}
