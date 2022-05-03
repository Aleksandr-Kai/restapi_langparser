package memstore

import (
	"restapi_langparser/internal/store"
)

type Store struct {
	ProxyRepository  store.IProxyRepository
	DomainRepository store.IDomainRepository
}

func New() *Store {
	return &Store{
		DomainRepository: NewDomainRepository(),
		ProxyRepository:  NewProxyRepository(),
	}
}

func (s *Store) Proxy() store.IProxyRepository {
	return s.ProxyRepository
}

func (s *Store) Domain() store.IDomainRepository {
	return s.DomainRepository
}

func (s *Store) Migrate() error {
	return nil
}
