package memstore

import (
	"errors"
	"fmt"
	"restapi_langparser/internal/model"
)

type DomainRepository struct {
	store  map[int]*model.Domain
	lastID uint
}

func NewDomainRepository() *DomainRepository {
	return &DomainRepository{
		store: make(map[int]*model.Domain),
	}
}

func (d *DomainRepository) getID() (uint, error) {
	prev := d.lastID
	d.lastID++
	cycle := prev < d.lastID
	for _, exists := d.store[int(d.lastID)]; exists; _, exists = d.store[int(d.lastID)] {
		if prev > d.lastID && cycle {
			return 0, errors.New("no id available")
		}
		d.lastID++
		cycle = prev < d.lastID
	}
	return d.lastID, nil
}

func (d *DomainRepository) Create(domains ...model.Domain) error {
	var err error
	for _, domain := range domains {
		domain.ID, err = d.getID()
		if err != nil {
			return fmt.Errorf(err.Error())
		}
		d.store[int(domain.ID)] = &domain
	}
	return nil
}

func (d *DomainRepository) Read(limit, offset int) ([]model.Domain, error) {
	if len(d.store) < offset+limit {
		return nil, errors.New("out of range")
	}
	res := make([]model.Domain, 0)
	for _, item := range d.store {
		if offset > 0 {
			offset--
			continue
		}
		res = append(res, *item)
		limit--
		if limit == 0 {
			break
		}
	}
	return res, nil
}

func (d *DomainRepository) Update(target model.Domain) error {
	return nil
}

func (d *DomainRepository) Delete(target ...model.Domain) error {
	for _, t := range target {
		delete(d.store, int(t.ID))
	}
	return nil
}

func (d *DomainRepository) FindByID(id int) (*model.Domain, error) {
	domain, exists := d.store[id]
	if !exists {
		return nil, fmt.Errorf("record with id=%d not found", id)
	}
	return domain, nil
}

func (d *DomainRepository) FindByTagLang(lang string) ([]model.Domain, error) {
	return nil, errors.New("method not implemented")
}

func (d *DomainRepository) FindBySMLang(lang string) ([]model.Domain, error) {
	return nil, errors.New("method not implemented")
}

func (d *DomainRepository) FindByContentLang(lang string) ([]model.Domain, error) {
	return nil, errors.New("method not implemented")
}

func (d *DomainRepository) FindByHost(hosts ...string) ([]model.Domain, error) {
	return nil, errors.New("method not implemented")
}

func (d *DomainRepository) GetUserRequest() (*model.Domain, error) {
	return nil, errors.New("method not implemented")
}

func (d *DomainRepository) GetErrorHost() (*model.Domain, error) {
	return nil, errors.New("method not implemented")
}

func (d *DomainRepository) GetListHost() (*model.Domain, error) {
	return nil, errors.New("method not implemented")
}

func (d *DomainRepository) CreateWithHost(from string, hosts ...string) ([]model.Domain, error) {
	return nil, errors.New("method not implemented")
}
