package sqlstore

import (
	"restapi_langparser/internal/model"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DomainRepository struct {
	db *gorm.DB
}

func NewDomainRepository(db *gorm.DB) *DomainRepository {
	return &DomainRepository{
		db: db,
	}
}

func (d *DomainRepository) Create(domains ...model.Domain) error {
	return d.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "host"}},
		DoNothing: true,
	}).Create(domains).Error
}

func (d *DomainRepository) Update(target model.Domain) error {
	tx := d.db.Omit("id")

	switch {
	case target.ID == 0 && target.Host != "":
		var id int64
		d.db.Table("domains").Select("id").Where("host=?", target.Host).Scan(&id)
		target.ID = uint(id)
	case target.ID > 0 && target.Host == "":
		tx = tx.Omit("host")
	}

	return tx.Save(&target).Error
}

func (d *DomainRepository) Read(limit, offset int) ([]model.Domain, error) {
	domains := make([]model.Domain, 0)
	tx := d.db
	err := tx.Model(&model.Domain{}).
		Preload(clause.Associations).
		Limit(limit).
		Offset(offset).
		Order("id").
		Find(&domains).
		Error
	if err != nil {
		return nil, err
	}

	return domains, nil
}

func (d *DomainRepository) Delete(target ...model.Domain) error {
	err := d.db.Delete(target).Error
	if err != nil {
		return err
	}
	d.db.Model(&model.Request{}).Select("domain_id").Joins("left join domains as d on d.id=r.domain_id where d.id isnull").Find(&model.Domain{})
	err = d.db.Exec(`delete from requests where domain_id in (select domain_id from requests as r left join domains as d on d.id=r.domain_id where d.id isnull)`).Error
	return err
}

func (d *DomainRepository) FindByID(id int) (*model.Domain, error) {
	domain := &model.Domain{}
	tx := d.db
	err := tx.Model(&model.Domain{}).
		Preload(clause.Associations).
		First(domain, id).
		Error
	if err != nil {
		return nil, err
	}

	return domain, nil
}

func (d *DomainRepository) getDomain(query string) (*model.Domain, error) {
	var domain *model.Domain
	tx := d.db

	err := tx.Where(query).
		First(domain).
		Error
	if err != nil {
		return nil, err
	}

	return domain, tx.Save(domain).Error
}

func (d *DomainRepository) GetUserRequest() (*model.Domain, error) {
	return d.getDomain("update < now() and from = 'user'")
}

func (d *DomainRepository) GetErrorHost() (*model.Domain, error) {
	return d.getDomain("update < now() and from = 'list' and response_code != 'ok' and response_code != ''")
}

func (d *DomainRepository) GetListHost() (*model.Domain, error) {
	return d.getDomain("update < now() and from = 'list' and response_code = ''")
}

func (d *DomainRepository) FindByTagLang(lang string) ([]model.Domain, error) {
	domains := make([]model.Domain, 0)
	tx := d.db
	err := tx.Model(&model.Domain{}).Where(`"tags_languages" like ?`, "%"+lang+"%").Find(&domains).Error
	if err != nil {
		return nil, err
	}
	return domains, nil
}

func (d *DomainRepository) FindBySMLang(lang string) ([]model.Domain, error) {
	domains := make([]model.Domain, 0)
	tx := d.db
	err := tx.Model(&model.Domain{}).Where(`"sitemap_languages" like ?`, "%"+lang+"%").Find(&domains).Error
	if err != nil {
		return nil, err
	}
	return domains, nil
}

func (d *DomainRepository) FindByContentLang(lang string) ([]model.Domain, error) {
	lang = strings.ToUpper(lang)
	domains := make([]model.Domain, 0)
	tx := d.db
	err := tx.Model(&model.Domain{}).Where(`"contentLang" = ?`, lang).Find(&domains).Error
	if err != nil {
		return nil, err
	}
	return domains, nil
}

func (d *DomainRepository) FindByHost(hosts ...string) ([]model.Domain, error) {
	var domains []model.Domain

	err := d.db.Preload(clause.Associations).Where("host in ?", hosts).Find(&domains).Error
	if err != nil {
		return nil, err
	}

	return domains, nil
}

func (d *DomainRepository) CreateWithHost(hosts ...string) ([]model.Domain, error) {
	domains := make([]model.Domain, 0)
	for _, host := range hosts {
		domains = append(domains, model.Domain{Host: host})
	}

	if err := d.Create(domains...); err != nil {
		return nil, err
	}

	return d.FindByHost(hosts...)
}
