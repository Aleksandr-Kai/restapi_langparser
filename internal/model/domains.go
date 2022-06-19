package model

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"gorm.io/gorm"
	"strings"
)

const (
	ResponseOk       = "ok"
	ResponseBan      = "ban"
	ResponseNotExist = "not exist"
	ResponseError    = "error"
	ResponseNull     = ""

	langSeparator = ","
)

type Domain struct {
	gorm.Model               `json:"-"`
	Host                     string    `json:"host" gorm:"column:host;unique"`
	ResponseCode             string    `json:"responseCode" gorm:"column:response_code"`
	ErrorCount               int       `json:"errorCount" gorm:"column:error_count"`
	ContentLanguage          string    `json:"contentLang" gorm:"column:content_lang"`
	TagsLanguages            []string  `json:"tagLanguages,omitempty" gorm:"-"`
	SitemapLanguages         []string  `json:"sitemapLanguages,omitempty" gorm:"-"`
	BlockerName              string    `json:"blockerName,omitempty" gorm:"column:blocker_name"`
	IP                       string    `json:"ip" gorm:"column:ip"`
	TagsLanguagesInternal    string    `json:"-" gorm:"column:tags_languages"`
	SitemapLanguagesInternal string    `json:"-" gorm:"column:sitemap_languages"`
	Requests                 []Request `json:"-" gorm:"foreignKey:DomainID;constraint:OnDelete:CASCADE;"`
	Queue                    Queue     `json:"-" gorm:"foreignKey:DomainID;constraint:OnDelete:CASCADE;"`
}

func (d *Domain) Validate() error {
	return validation.ValidateStruct(
		d,
		validation.Field(&d.ID, validation.Required, validation.Min(0)),
		validation.Field(&d.Host, validation.Required, is.URL),
		validation.Field(&d.IP, is.IP),
	)
}

func (d *Domain) languagesToString() {
	d.TagsLanguagesInternal = strings.ToUpper(strings.Join(d.TagsLanguages, langSeparator))
	d.SitemapLanguagesInternal = strings.ToUpper(strings.Join(d.SitemapLanguages, langSeparator))
}

func (d *Domain) languagesToSlice() {
	if d.TagsLanguagesInternal != "" {
		d.TagsLanguages = strings.Split(d.TagsLanguagesInternal, langSeparator)
	}
	if d.SitemapLanguagesInternal != "" {
		d.SitemapLanguages = strings.Split(d.SitemapLanguagesInternal, langSeparator)
	}
}

func (d *Domain) BeforeCreate(*gorm.DB) (err error) {
	if d.Host == "" {
		return errors.New("empty host")
	}
	return nil
}

func (d *Domain) BeforeUpdate(*gorm.DB) (err error) {
	err = validation.ValidateStruct(
		d,
		validation.Field(&d.ContentLanguage, validation.RuneLength(0, 2)),
		validation.Field(&d.TagsLanguages, validation.Each(validation.RuneLength(0, 2))),
		validation.Field(&d.SitemapLanguages, validation.Each(validation.RuneLength(0, 2))),
	)
	if err != nil {
		return err
	}

	d.languagesToString()

	return nil
}

func (d *Domain) AfterFind(*gorm.DB) (err error) {
	d.languagesToSlice()
	return nil
}

func CreateDomainsList(hosts []string) []Domain {
	domains := make([]Domain, len(hosts))
	for i, host := range hosts {
		domains[i].Host = host
	}
	return domains
}
