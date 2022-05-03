package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type Domain struct {
	ID                  int64          `json:"id" gorm:"primaryKey;column:id"`
	URL                 string         `json:"url" gorm:"column:url;unique"`
	ResponseCode        uint8          `json:"response_code" gorm:"column:response_code"`
	ErrorCount          uint           `json:"error_count" gorm:"column:error_count"`
	ContentLanguage     string         `json:"content_lang" gorm:"column:content_lang"`
	TabTagsLanguages    []TagsLangs    `json:"-" gorm:"foreignKey:DomainID"`
	TagsLanguages       []string       `json:"tag_languages" gorm:"-"`
	TabSitemapLanguages []SitemapLangs `json:"-" gorm:"foreignKey:DomainID"`
	SitemapLanguages    []string       `json:"sitemap_languages" gorm:"-"`
	BlockerName         string         `json:"blocker_name" gorm:"column:blocker_name"`
	IP                  string         `json:"ip" gorm:"column:ip"`
}

func (d *Domain) Validate() error {
	return validation.ValidateStruct(
		d,
		validation.Field(&d.ID, validation.Required, validation.Min(0)),
		validation.Field(&d.URL, validation.Required, is.URL),
		validation.Field(&d.IP, is.IP),
	)
}

func (d *Domain) LanguagesAsStrings() {
	d.TagsLanguages = make([]string, len(d.TabSitemapLanguages))
	for ind := range d.TabSitemapLanguages {
		d.TagsLanguages[ind] = d.TabSitemapLanguages[ind].Lang
	}

	d.SitemapLanguages = make([]string, len(d.TabSitemapLanguages))
	for ind := range d.TabSitemapLanguages {
		d.SitemapLanguages[ind] = d.TabSitemapLanguages[ind].Lang
	}
}
