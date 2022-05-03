package store

import (
	"restapi_langparser/internal/model"
)

type LangSource int

type IProxyRepository interface {
	Add(list []model.Proxy) error
	List(limit, offset int) ([]model.Proxy, error)
	Remove(ids []int) error
	Update(list []model.Proxy) error
}

type IDomainRepository interface {
	Add(domain model.Domain) (int64, error)
	List(limit, offset int64) ([]model.Domain, error)
	Remove(ids []int64) error
	Update(list []model.Domain) error

	FindByID(id int64) (*model.Domain, error)
	FindByLang(lang string, from LangSource) ([]model.Domain, error)
	FindByHost(host string) (*model.Domain, error)
}

type IStore interface {
	Migrate() error
	Proxy() IProxyRepository
	Domain() IDomainRepository
}
