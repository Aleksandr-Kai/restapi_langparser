package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

const (
	HTTPS  ProxyType = "https"
	Socks4 ProxyType = "socks4"
	Socks5 ProxyType = "socks5"
)

type ProxyType string

type Proxy struct {
	ID   int       `json:"id"`
	URL  string    `json:"url"`
	Type ProxyType `json:"type"`
}

func (p *Proxy) Validate() error {
	return validation.ValidateStruct(
		p,
		validation.Field(&p.URL, validation.Required, is.URL),
	)
}
