package proxyprovider

import (
	"restapi_langparser/internal/config"
	"restapi_langparser/internal/model"
)

type ProxyProvider struct {
	config    *config.Config
	proxyList map[string]uint
}

func New(config *config.Config) *ProxyProvider {
	return &ProxyProvider{
		config: config,
	}
}

func (p *ProxyProvider) Get() *model.Proxy {
	return nil
}
