package langfinder

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/temoto/robotstxt"
	"io/ioutil"
	"net/http"
	"net/url"
	"restapi_langparser/internal/config"
	"restapi_langparser/internal/model"
	"restapi_langparser/internal/parser"
	"restapi_langparser/internal/proxyprovider"
	"restapi_langparser/internal/store"
	"sync"
	"time"
)

const (
	maxErrorsCount = 10
)

type LangFinder struct {
	store         store.IStore
	config        *config.Config
	threadLimit   chan interface{}
	proxyProvider *proxyprovider.ProxyProvider
	callbacks     struct {
		sync.RWMutex
		m map[string]string
	}
}

func New(store store.IStore, config *config.Config) *LangFinder {
	return &LangFinder{
		store:         store,
		config:        config,
		proxyProvider: proxyprovider.New(config, store),
		callbacks: struct {
			sync.RWMutex
			m map[string]string
		}{m: make(map[string]string)},
	}
}

func (f *LangFinder) NewTask(callback string, hosts ...string) (string, error) {
	if len(hosts) == 0 {
		return "", errors.New("empty host list")
	}
	domains := make([]model.Domain, len(hosts))
	for i, host := range hosts {
		domains[i] = model.Domain{
			Host: host,
		}
	}

	err := f.store.AddDomains(&domains)

	if err != nil {
		return "", err
	}

	//if queuePrior == model.RequestQueue {
	//	domains, err = f.store.Domain().FindByHost(hosts...)
	//	if err != nil {
	//		return "", err
	//	}
	//	requestCode, err := f.store.Request().Create(domains)
	//	if err != nil {
	//		return "", err
	//	}
	//
	//	if callback != "" {
	//		f.callbacks.Lock()
	//		f.callbacks.m[requestCode] = callback
	//		f.callbacks.Unlock()
	//	}
	//	return requestCode, nil
	//}

	return "", nil
}

func (f *LangFinder) getSitemapURLs(uRL string) ([]string, error) {
	u, _ := url.Parse(uRL)
	uRL = fmt.Sprintf("%s://%s/robots.txt", u.Scheme, u.Host)

	client := http.Client{
		Timeout: time.Second * 5,
	}

	req, err := http.NewRequest(http.MethodGet, uRL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	robots, err := robotstxt.FromResponse(resp)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	if len(robots.Sitemaps) == 0 {
		return []string{fmt.Sprintf("%s://%s/sitemap.xml", u.Scheme, u.Host)}, nil
	}

	return robots.Sitemaps, nil
}

func (f *LangFinder) taskWorker(domain model.Domain, taskLimiter chan interface{}) {
	defer func() {
		<-taskLimiter
	}()

	req, err := http.NewRequest(http.MethodGet, domain.Host, nil)
	if err != nil {
		logrus.Errorf("Creqte http request fail: %s", err)
		return
	}

	proxy := f.proxyProvider.Get()
	if proxy == nil {
		return
	}
	defer func() {
		f.proxyProvider.Release(proxy)
	}()

	client := createClient(proxy, f.config.ResponseTimeout)

	defer func() {
		err = f.store.Domain().Update(domain)
		if err != nil {
			logrus.Errorf("DB update fail: %s", err)
		}
	}()

	// request page
	resp, err := client.Do(req)
	if err != nil {
		logrus.Errorf("URL visit fail: %s", err)
		if domain.ResponseCode == model.ResponseError {
			domain.ErrorCount++
		} else {
			domain.ErrorCount = 1
		}
		return
	}

	switch resp.StatusCode {
	case http.StatusOK:
		domain.ResponseCode = model.ResponseOk
	case http.StatusNotFound:
		domain.ResponseCode = model.ResponseNotExist
	default:
		domain.ResponseCode = model.ResponseError
	}

	page, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	domain.ContentLanguage, err = parser.GetContentLang(bytes.NewReader(page))
	if err != nil {
		domain.ErrorCount++
		domain.ResponseCode = model.ResponseError
	}

	domain.TagsLanguages, err = parser.GetLangsInTags(bytes.NewReader(page))
	if err != nil {
		domain.ErrorCount++
		domain.ResponseCode = model.ResponseError
	}

	err = f.store.Domain().Update(domain)
	if err != nil {
		logrus.Error(err)
	}
}

func (f *LangFinder) taskManager(limit int) {
	errorsCnt := 0
	//updateLimit := int(float64(limit) * 0.6)
	taskLimiter := make(chan interface{}, limit-1)
	for {
		if errorsCnt > maxErrorsCount {
			logrus.Errorf("error limit exceeded")
			return
		}
		queueDomain, err := f.store.Domain().GetUserRequest()
		if err != nil {
			logrus.Errorf("Get queue fail: %s", err)
			errorsCnt++
			continue
		}

		if queueDomain == nil { // user request not found
			queueDomain, err = f.store.Domain().GetErrorHost()
			if err != nil {
				logrus.Errorf("Get queue fail: %s", err)
				errorsCnt++
				continue
			}
		}

		taskLimiter <- 0
		go f.taskWorker(*queueDomain, taskLimiter)

		errorsCnt = 0

		//if task == nil {
		//	time.Sleep(time.Second)
		//	continue
		//}
		//qt, ok := task.(taskItem)
		//if !ok {
		//	continue
		//}
		//
		//taskLimiter <- 0
		//go f.taskWorker(taskLimiter)
	}
	//resp, err := f.requester.GetPageData(uRL)
	//if err != nil {
	//	log.Println("parse error:", err)
	//	return
	//}
	//contentLang, err := parser.GetContentLang(resp.Page())
	//if err != nil {
	//	log.Println("parse error:", err)
	//	return
	//}
	//tagLangs, err := parser.GetLangsInTags(resp.Page())
	//if err != nil {
	//	log.Println("parse error:", err)
	//	return
	//}
	//
	//dom := model.Domain{
	//	URL:             uRL,
	//	ResponseCode:    uint16(resp.ResponseCode),
	//	ContentLanguage: contentLang,
	//}
	//
	//// add langs form header and meta
	//dom.TabTagsLanguages = make([]model.TagsLangs, len(tagLangs))
	//for i, tl := range tagLangs {
	//	dom.TabTagsLanguages[i] = model.TagsLangs{
	//		Lang: tl,
	//	}
	//}
	//for _, tl := range resp.HeaderLang {
	//	dom.TabTagsLanguages = append(dom.TabTagsLanguages, model.TagsLangs{
	//		Lang: tl,
	//	})
	//}
	//
	//// add langs from sitemap
	//
	//id, err := f.store.Domain().Add(dom)
	//
	//if err != nil {
	//	log.Println("db error:", err)
	//	return
	//}
	//
	//f.userRequestResults[requestCode].Add(id)
}
