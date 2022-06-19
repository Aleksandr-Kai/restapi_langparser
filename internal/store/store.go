package store

import (
	"restapi_langparser/internal/model"
)

type IProxyRepository interface {
	Create(list []model.Proxy) error
	Read(limit, offset int) ([]model.Proxy, error)
	Update(list []model.Proxy) error
	Delete(ids []int) error
}

type IDomainRepository interface {
	Create(domains ...model.Domain) error
	Read(limit, offset int) ([]model.Domain, error)
	Update(target model.Domain) error
	Delete(target ...model.Domain) error

	FindByID(id int) (*model.Domain, error)
	FindByTagLang(lang string) ([]model.Domain, error)
	FindBySMLang(lang string) ([]model.Domain, error)
	FindByContentLang(lang string) ([]model.Domain, error)
	FindByHost(hosts ...string) ([]model.Domain, error)
	GetUserRequest() (*model.Domain, error)
	GetErrorHost() (*model.Domain, error)
	GetListHost() (*model.Domain, error)
	CreateWithHost(from string, hosts ...string) ([]model.Domain, error)
}

type IStore interface {
	Migrate() error

	Proxy() IProxyRepository
	Domain() IDomainRepository

	AddDomains(list *[]model.Domain) error
	GetDomains(hosts []string) ([]model.Domain, error)
	GetFromQueue() *model.Domain
	SaveDomain(domain model.Domain) error
	CreateRequest(list []model.Domain, callback *string) (requestCode string, err error)
	GetRequest(requestCode string) ([]model.Domain, error)
}
