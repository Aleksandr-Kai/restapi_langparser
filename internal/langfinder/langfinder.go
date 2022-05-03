package langfinder

import (
	"crypto/md5"
	"errors"
	"fmt"
	"log"
	"restapi_langparser/internal/config"
	"restapi_langparser/internal/model"
	"restapi_langparser/internal/proxyprovider"
	"restapi_langparser/internal/store"
	"strings"
	"time"
)

type LangFinder struct {
	store         store.IStore
	config        *config.Config
	proxyProvider *proxyprovider.ProxyProvider
	threadLimit   chan interface{}

	results map[string]*requestResult
}

func New(store store.IStore, config *config.Config) *LangFinder {
	lf := &LangFinder{
		store:         store,
		config:        config,
		proxyProvider: proxyprovider.New(config),
		threadLimit:   make(chan interface{}, config.MaxThreads-1),
		results:       make(map[string]*requestResult),
	}
	go lf.schedule()
	return lf
}

func (f *LangFinder) schedule() {
	for now := range time.Tick(time.Hour) {
		codeList := make([]string, 0)
		for code, res := range f.results {
			if res.expiration.Before(now) {
				codeList = append(codeList, code)
			}
		}
		for _, code := range codeList {
			delete(f.results, code)
		}
	}
}

func (f *LangFinder) getHash(data string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(data)))
}

func (f *LangFinder) NewTask(callback string, urls ...string) string {
	requestCode := f.getHash(callback + strings.Join(urls, ""))
	log.Println("create new task", requestCode, urls)
	go func() {
		f.results[requestCode] = NewResult(len(urls))
		for _, url := range urls {
			f.threadLimit <- 0
			log.Println("create new thread", url)
			proxy := f.proxyProvider.Get()
			go f.worker(url, proxy, requestCode)
		}
	}()

	return requestCode
}

func (f *LangFinder) worker(url string, proxy *model.Proxy, requestCode string) {
	defer func() {
		<-f.threadLimit
	}()
	log.Println("start worker", url, requestCode)
	id, err := f.store.Domain().Add(model.Domain{
		URL: url,
	})
	if err != nil {
		log.Println("db error:", err)
		return
	}

	f.results[requestCode].Add(id)
}

func (f *LangFinder) GetResult(requestCode string) ([]model.Domain, error) {
	result, exists := f.results[requestCode]
	if !exists {
		return nil, fmt.Errorf("result with code %v not found", requestCode)
	}
	if !result.Ready() {
		return nil, errors.New("result not ready")
	}
	ids := f.results[requestCode].DomainIDs
	res := make([]model.Domain, len(ids))
	for i, id := range ids {
		domain, err := f.store.Domain().FindByID(id)
		if err != nil {
			log.Println("find by id fail:", err)
			continue
		}
		res[i] = *domain
	}
	return res, nil
}
