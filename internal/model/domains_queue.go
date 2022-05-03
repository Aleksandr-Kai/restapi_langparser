package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type Queue struct {
	URL string `json:"URL"`
}

func (q *Queue) Validate() error {
	return validation.ValidateStruct(
		q,
		validation.Field(&q.URL, validation.Required, is.URL),
	)
}
