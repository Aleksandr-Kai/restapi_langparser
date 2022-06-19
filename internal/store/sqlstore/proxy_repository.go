package sqlstore

import (
	"errors"
	"gorm.io/gorm"
	"restapi_langparser/internal/model"
)

type ProxyRepository struct {
	db *gorm.DB
}

func NewProxyRepository(db *gorm.DB) *ProxyRepository {
	return &ProxyRepository{
		db: db,
	}
}

func (p *ProxyRepository) Create(list []model.Proxy) error {
	batch := make([]model.Proxy, 0, len(list))
	for _, proxy := range list {
		if proxy.Validate() == nil { // skip if not valid
			batch = append(batch, proxy)
		}
	}
	return p.db.Create(batch).Error
}

func (p *ProxyRepository) Read(limit, offset int) ([]model.Proxy, error) {
	return nil, errors.New("not implemented")
}

func (p *ProxyRepository) Delete(ids []int) error {
	return errors.New("not implemented")
}

func (p *ProxyRepository) Update(list []model.Proxy) error {
	return errors.New("not implemented")
}

func (p *ProxyRepository) FindByID(id int64) (*model.Proxy, error) {
	proxy := &model.Proxy{}
	tx := p.db
	err := tx.Model(&model.Proxy{}).
		First(proxy, id).
		Error
	if err != nil {
		return nil, err
	}
	return proxy, nil
}
