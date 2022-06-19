package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"strings"
)

const (
	HTTPS   = "https"
	Socks4  = "socks4"
	Socks5  = "socks5"
	NoProxy = "noproxy"
)

type Proxy struct {
	ID       int    `gorm:"primaryKey;column:id" json:"-"`
	IP       string `gorm:"column:ip" json:"ip"`
	Port     string `gorm:"column:port" json:"port"`
	Login    string `gorm:"column:login" json:"login"`
	Password string `gorm:"column:password" json:"password"`
	Scheme   string `gorm:"column:type" json:"type"`
}

func (p *Proxy) Validate() error {
	return validation.ValidateStruct(
		p,
		validation.Field(p.IP, validation.Required, is.IPv4),
		validation.Field(p.Port, validation.Required, is.Digit),
		validation.Field(p.Type, validation.Required),
	)
}

func (p *Proxy) Type() string {
	return strings.ToLower(p.Scheme)
}
