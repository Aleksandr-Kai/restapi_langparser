package proxyprovider

import (
	"restapi_langparser/internal/config"
	"restapi_langparser/internal/model"
	"restapi_langparser/internal/store"
)

type ProxyProvider struct {
	config       *config.Config
	store        store.IStore
	activeProxy  []proxyItem
	noProxyCount uint
}

type proxyItem struct {
	proxy       model.Proxy
	threadCount uint
}

func New(config *config.Config, store store.IStore) *ProxyProvider {
	pp := &ProxyProvider{
		config:      config,
		store:       store,
		activeProxy: make([]proxyItem, 0),
	}
	pp.updateList()
	return pp
}

func (p *ProxyProvider) Get() *model.Proxy {
	for _, item := range p.activeProxy {
		if item.threadCount < p.config.ThreadsPerProxy {
			item.threadCount++
			return &item.proxy
		}
	}

	if p.noProxyCount < p.config.ThreadsPerProxy {
		p.noProxyCount++
		return &model.Proxy{
			ID:     -1,
			Scheme: model.NoProxy,
		}
	}

	return nil
}

func (p *ProxyProvider) Release(proxy *model.Proxy) {
	if proxy == nil {
		return
	}
	if proxy.Type() == model.NoProxy {
		if p.noProxyCount > 0 {
			p.noProxyCount--
		}
		return
	}
	for _, item := range p.activeProxy {
		if item.proxy.ID == proxy.ID {
			if item.threadCount > 0 {
				item.threadCount--
			}
			return
		}
	}
}

func (p *ProxyProvider) updateList() {
	list, err := p.store.Proxy().Read(0, 0)
	if err != nil {
		return
	}
	newList := make([]proxyItem, len(list))
	for i, proxy := range list {
		for _, item := range p.activeProxy {
			if proxy.ID == item.proxy.ID {
				newList[i].threadCount = item.threadCount
				break
			}
		}
		newList[i].proxy = proxy
	}

	p.activeProxy = newList
}
