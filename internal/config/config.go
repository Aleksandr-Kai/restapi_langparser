package config

import (
	"time"
)

const (
	MemStore StoreType = "memstore"
	SQLStore StoreType = "sqlstore"

	defaultThreadsPerProxy           = 10
	defaultMaxThreads                = 100
	defaultMaxRequestToBadDomain     = 60 // %
	defaultResponseTimeout           = time.Second * 15
	defaultDeadProxyRefresh          = time.Minute * 70
	defaultDomainWithErrorRefresh    = time.Hour * 72
	defaultDomainWithDNSErrorRefresh = time.Hour * 336
	defaultRebannedDomainRefresh     = time.Hour * 72
)

type StoreType string

type Config struct {
	Type        StoreType
	BindAddr    string `toml:"bind_addr"`
	DatabaseURL string `toml:"database_url"`

	UseIP                     bool
	ThreadsPerProxy           uint
	MaxThreads                uint
	MaxRequestToBadDomain     uint
	ResponseTimeout           time.Duration
	DeadProxyRefresh          time.Duration
	DomainWithErrorRefresh    time.Duration
	DomainWithDNSErrorRefresh time.Duration
	RebannedDomainRefresh     time.Duration
}

func New() *Config {
	return &Config{
		Type:                      SQLStore,
		BindAddr:                  ":8080",
		UseIP:                     false,
		ThreadsPerProxy:           defaultThreadsPerProxy,
		MaxThreads:                defaultMaxThreads,
		MaxRequestToBadDomain:     defaultMaxRequestToBadDomain,
		ResponseTimeout:           defaultResponseTimeout,
		DeadProxyRefresh:          defaultDeadProxyRefresh,
		DomainWithErrorRefresh:    defaultDomainWithErrorRefresh,
		DomainWithDNSErrorRefresh: defaultDomainWithDNSErrorRefresh,
		RebannedDomainRefresh:     defaultRebannedDomainRefresh,
	}
}
