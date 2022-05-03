package sqlstore

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"reflect"
	"restapi_langparser/internal/model"
	"restapi_langparser/internal/store"
	"strings"
)

const (
	FromContent = iota
	FromTags
	FromSitemap
)

type DomainRepository struct {
	db *gorm.DB
}

func NewDomainRepository(db *gorm.DB) store.IDomainRepository {
	return &DomainRepository{
		db: db,
	}
}

func (d *DomainRepository) Add(domain model.Domain) (int64, error) {
	tx := d.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{
				Name: "url",
			},
		},
		UpdateAll: true,
	}).Create(&domain)
	err := tx.Error
	tx.Save(domain)
	return domain.ID, err
}

func (d *DomainRepository) List(limit, offset int64) ([]model.Domain, error) {
	domains := make([]model.Domain, 0)
	tx := d.db
	err := tx.Model(&model.Domain{}).
		Preload(clause.Associations).
		Limit(int(limit)).
		Offset(int(offset)).
		Find(&domains).
		Group("id").
		Error
	if err != nil {
		return nil, err
	}

	for i := range domains {
		domains[i].LanguagesAsStrings()
	}

	return domains, nil
}

func (d *DomainRepository) Remove(ids []int64) error {
	return nil
}

func (d *DomainRepository) Update(list []model.Domain) error {
	return nil
}

func (d *DomainRepository) FindByID(id int64) (*model.Domain, error) {
	domain := &model.Domain{}
	tx := d.db
	err := tx.Model(&model.Domain{}).
		Preload(clause.Associations).
		First(domain, id).
		Error
	if err != nil {
		return nil, err
	}

	domain.LanguagesAsStrings()

	return domain, nil
}

func (d *DomainRepository) FindByLang(lang string, from store.LangSource) ([]model.Domain, error) {
	lang = strings.ToUpper(lang)
	domains := make([]model.Domain, 0)
	tx := d.db

	var err error
	if from == FromContent {
		err = tx.Model(&model.Domain{}).
			Preload(clause.Associations).
			Find(&domains, model.Domain{ContentLanguage: lang}).
			Error
	} else {
		var field reflect.StructField
		switch from {
		case FromTags:
			field, _ = reflect.TypeOf(model.Domain{}).Elem().FieldByName("TabTagsLanguages")
		case FromSitemap:
			field, _ = reflect.TypeOf(model.Domain{}).Elem().FieldByName("TabSitemapLanguages")
		}
		err = tx.Model(&model.Domain{}).
			Joins(fmt.Sprintf("inner join %s as l on l.domain_id=domains.id and l.lang=?", string(field.Tag)), lang).
			Preload(clause.Associations).
			Find(&domains).
			Error
	}

	if err != nil {
		return nil, err
	}
	for i := range domains {
		domains[i].LanguagesAsStrings()
	}
	return domains, nil
}

func (d *DomainRepository) FindByHost(host string) (*model.Domain, error) {
	domain := &model.Domain{URL: host}
	tx := d.db

	err := tx.Preload(clause.Associations).First(domain, domain).Error
	if err != nil {
		return nil, err
	}
	domain.LanguagesAsStrings()
	return domain, nil
}
