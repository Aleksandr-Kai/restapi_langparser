package sqlstore

import (
	"errors"
	"restapi_langparser/internal/model"
	"restapi_langparser/internal/store"

	"gorm.io/gorm"
)

type ProxyRepository struct {
	db *gorm.DB
}

func NewProxyRepository(db *gorm.DB) store.IProxyRepository {
	return &ProxyRepository{
		db: db,
	}
}

func (p *ProxyRepository) Add(list []model.Proxy) error {
	batch := make([]model.Proxy, 0, len(list))
	for _, proxy := range list {
		if proxy.Validate() == nil { // skip if not valid
			batch = append(batch, proxy)
		}
	}
	return p.db.Create(batch).Error
}

func (p *ProxyRepository) List(limit, offset int) ([]model.Proxy, error) {
	return nil, errors.New("method not implemented")
}

func (p *ProxyRepository) Remove(ids []int) error {
	return nil
}

func (p *ProxyRepository) Update(list []model.Proxy) error {
	return nil
}
