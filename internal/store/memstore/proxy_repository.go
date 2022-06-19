package memstore

import (
	"errors"
	"fmt"
	"restapi_langparser/internal/model"
	"restapi_langparser/internal/store"
)

type ProxyRepository struct {
	store  map[int]model.Proxy
	lastID int
}

func (p *ProxyRepository) getID() (int, error) {
	prev := p.lastID
	p.lastID++
	cycle := prev < p.lastID
	for _, exists := p.store[p.lastID]; exists; _, exists = p.store[p.lastID] {
		if prev > p.lastID && cycle {
			return 0, errors.New("no id available")
		}
		p.lastID++
		cycle = prev < p.lastID
	}
	return p.lastID, nil
}

func NewProxyRepository() store.IProxyRepository {
	return &ProxyRepository{
		store: make(map[int]model.Proxy),
	}
}

func (p *ProxyRepository) Add(list []model.Proxy) error {
	var err error
	for i, proxy := range list {
		proxy.ID, err = p.getID()
		if err != nil {
			return fmt.Errorf("%v %d records processed", err, i-1)
		}
		p.store[proxy.ID] = proxy
	}

	return nil
}

func (p *ProxyRepository) List(limit, offset int) ([]model.Proxy, error) {
	if len(p.store) < offset+limit {
		return nil, errors.New("out of range")
	}
	res := make([]model.Proxy, 0)
	for _, item := range p.store {
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

func (p *ProxyRepository) Remove(ids []int) error {
	return nil
}

func (p *ProxyRepository) Update(list []model.Proxy) error {
	return nil
}

func (p *ProxyRepository) FindByID(id int64) (*model.Proxy, error) {
	domain, exists := p.store[int(id)]
	if !exists {
		return nil, fmt.Errorf("record with id=%d not found", id)
	}
	return &domain, nil
}
