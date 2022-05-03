package sqlstore

import (
	"restapi_langparser/internal/model"
	"restapi_langparser/internal/store"

	_ "github.com/lib/pq" //nolint:goimports
	"gorm.io/gorm"
)

type Store struct {
	db               *gorm.DB
	ProxyRepository  store.IProxyRepository
	DomainRepository store.IDomainRepository
}

func New(db *gorm.DB) *Store {
	return &Store{
		db:               db,
		DomainRepository: NewDomainRepository(db),
		ProxyRepository:  NewProxyRepository(db),
	}
}

func (s *Store) Proxy() store.IProxyRepository {
	return s.ProxyRepository
}

func (s *Store) Domain() store.IDomainRepository {
	return s.DomainRepository
}

func (s *Store) Migrate() error {
	m := s.db.Migrator()
	return m.AutoMigrate(&model.Proxy{}, &model.Domain{}, &model.Blocked{})
}
