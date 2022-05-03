package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type Blocked struct {
	URL string `json:"URL"`
}

func (b *Blocked) Validate() error {
	return validation.ValidateStruct(
		b,
		validation.Field(&b.URL, validation.Required, is.URL),
	)
}
