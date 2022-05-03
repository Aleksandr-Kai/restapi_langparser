package memstore

import (
	"errors"
	"fmt"
	"restapi_langparser/internal/model"
	"restapi_langparser/internal/store"
)

type DomainRepository struct {
	store  map[int64]model.Domain
	lastID int64
}

func NewDomainRepository() store.IDomainRepository {
	return &DomainRepository{
		store: make(map[int64]model.Domain),
	}
}

func (d *DomainRepository) getID() (int64, error) {
	prev := d.lastID
	d.lastID++
	cycle := prev < d.lastID
	for _, exists := d.store[d.lastID]; exists; _, exists = d.store[d.lastID] {
		if prev > d.lastID && cycle {
			return 0, errors.New("no id available")
		}
		d.lastID++
		cycle = prev < d.lastID
	}
	return d.lastID, nil
}

func (d *DomainRepository) Add(domain model.Domain) (int64, error) {
	var err error
	domain.ID, err = d.getID()
	if err != nil {
		return 0, fmt.Errorf(err.Error())
	}
	d.store[domain.ID] = domain
	return domain.ID, nil
}

func (d *DomainRepository) List(limit, offset int64) ([]model.Domain, error) {
	if int64(len(d.store)) < offset+limit {
		return nil, errors.New("out of range")
	}
	res := make([]model.Domain, 0)
	for _, item := range d.store {
		if offset > 0 {
			offset--
			continue
		}
		res = append(res, item)
		limit--
		if limit == 0 {
			break
		}
	}
	return res, nil
}

func (d *DomainRepository) Remove(ids []int64) error {
	for _, id := range ids {
		delete(d.store, id)
	}
	return nil
}

func (d *DomainRepository) Update(list []model.Domain) error {
	return nil
}

func (d *DomainRepository) FindByID(id int64) (*model.Domain, error) {
	domain, exists := d.store[id]
	if !exists {
		return nil, fmt.Errorf("record with id=%d not found", id)
	}
	return &domain, nil
}

func (d *DomainRepository) FindByLang(lang string, from store.LangSource) ([]model.Domain, error) {
	return nil, errors.New("method not implemented")
}

func (d *DomainRepository) FindByHost(host string) (*model.Domain, error) {
	return nil, errors.New("method not implemented")
}
