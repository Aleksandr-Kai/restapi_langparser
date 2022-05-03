package apiserver

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"restapi_langparser/internal/apistructs"
	"restapi_langparser/internal/config"
	"restapi_langparser/internal/langfinder"
	"restapi_langparser/internal/model"
	"restapi_langparser/internal/store"
	"strconv"
)

type ctxKey int8

type server struct {
	router *gin.Engine
	store  store.IStore
	finder *langfinder.LangFinder
}

func newServer(store store.IStore, config *config.Config) *server {
	s := &server{
		router: gin.New(),
		store:  store,
		finder: langfinder.New(store, config),
	}
	s.configureRouter()
	return s
}

func (s *server) configureRouter() {
	s.router.Use(gin.Recovery())
	s.router.POST("/domains", s.handleAddDomains)
	s.router.GET("/domains", s.handleGetDomains)
	s.router.GET("/domains/:id", s.handleGetDomain)
	s.router.GET("/results/:code", s.handleGetResult)
	s.router.PATCH("/domains/:id", s.handleUpdateDomain)
	s.router.DELETE("/domains/:id", s.handleDeleteDomain)

	s.router.POST("/proxy", s.handleAddProxy)
	s.router.GET("/proxy", s.handleGetProxyList)
	s.router.PATCH("/proxy/:id", s.handleUpdateProxy)
	s.router.DELETE("/proxy/:id", s.handleDeleteProxy)
}

func (s *server) Start(addr string) error {
	return s.router.Run(addr)
}

/*
	Handlers
*/

func (s *server) handleGetResult(c *gin.Context) {
	resp := &apistructs.APIResponse{}
	status := http.StatusInternalServerError
	defer func() {
		c.String(status, resp.String())
	}()

	status = http.StatusOK

	res, err := s.finder.GetResult(c.Param("code"))
	if err != nil {
		resp.CreateErrorf(err.Error())
		return
	}
	resp.Results = &apistructs.APIResults{
		Domains: res,
	}
}

func (s *server) handleGetDomains(c *gin.Context) {
	logrus.Infof("Request [domain list] from %v", c.ClientIP())
	resp := &apistructs.APIResponse{}
	status := http.StatusInternalServerError
	defer func() {
		c.String(status, resp.String())
	}()

	status = http.StatusOK
	domains, err := s.store.Domain().List(0, 0)
	if err != nil {
		resp.CreateErrorf(err.Error())
		return
	}
	resp.Results = &apistructs.APIResults{
		Domains: domains,
	}
}

func (s *server) handleGetDomain(c *gin.Context) {
	resp := &apistructs.APIResponse{}
	status := http.StatusInternalServerError
	defer func() {
		c.String(status, resp.String())
	}()

	status = http.StatusOK
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		resp.CreateErrorf("Incorrect ID (%s)", err.Error())
		return
	}
	domain, err := s.store.Domain().FindByID(int64(id))
	if err != nil {
		resp.CreateErrorf(err.Error())
		return
	}
	resp.Results = &apistructs.APIResults{
		Domains: []model.Domain{*domain},
	}
}

func (s *server) handleUpdateDomain(c *gin.Context) {
	resp := &apistructs.APIResponse{}
	status := http.StatusInternalServerError
	defer func() {
		c.String(status, resp.String())
	}()

	status = http.StatusOK
	resp.CreateMessagef("request %v", c.Param("id"))
}

func (s *server) handleDeleteDomain(c *gin.Context) {
	resp := &apistructs.APIResponse{}
	status := http.StatusInternalServerError
	defer func() {
		c.String(status, resp.String())
	}()

	status = http.StatusOK
	resp.CreateMessagef("request %v", c.Param("id"))
}

func (s *server) handleAddDomains(c *gin.Context) {
	var req apistructs.APIRequest
	resp := &apistructs.APIResponse{}
	status := http.StatusInternalServerError
	defer func() {
		c.String(status, resp.String())
	}()

	if err := c.BindJSON(&req); err != nil {
		resp.CreateErrorf(err.Error())
		status = http.StatusInternalServerError
		return
	}

	status = http.StatusOK
	resp.CreateMessagef("request code")
	//todo make request and parse languages
	resp.Results = &apistructs.APIResults{
		RequestCode: s.finder.NewTask(req.Callback, req.URLs...),
	}
}

func (s *server) handleAddProxy(c *gin.Context) {
	var request apistructs.APIRequest
	resp := &apistructs.APIResponse{}
	status := http.StatusInternalServerError
	defer func() {
		c.String(status, resp.String())
	}()

	if err := c.BindJSON(&request); err != nil {
		resp.CreateErrorf(err.Error())
		status = http.StatusInternalServerError
		return
	}

	proxyList := make([]model.Proxy, len(request.URLs))
	for i, url := range request.URLs {
		proxyList[i] = model.Proxy{
			URL: url,
		}
	}
	err := s.store.Proxy().Add(proxyList)
	if err != nil {
		resp.CreateErrorf(err.Error())
		return
	}

	lst, err := s.store.Proxy().List(0, 0)
	if err != nil {
		resp.CreateErrorf(err.Error())
		return
	}
	resp.CreateMessagef("proxy %v", lst)
	status = http.StatusOK
}

func (s *server) handleGetProxyList(c *gin.Context) {
	resp := &apistructs.APIResponse{}
	status := http.StatusInternalServerError
	defer func() {
		c.String(status, resp.String())
	}()

	lst, err := s.store.Proxy().List(0, 0)
	if err != nil {
		resp.CreateErrorf(err.Error())
		return
	}
	resp.Results = &apistructs.APIResults{
		Proxy: lst,
	}
	resp.CreateMessagef("proxy %v", lst)

	status = http.StatusOK
}

func (s *server) handleUpdateProxy(c *gin.Context) {
	resp := &apistructs.APIResponse{}
	status := http.StatusInternalServerError
	defer func() {
		c.String(status, resp.String())
	}()

	status = http.StatusOK
}

func (s *server) handleDeleteProxy(c *gin.Context) {
	resp := &apistructs.APIResponse{}
	status := http.StatusInternalServerError
	defer func() {
		c.String(status, resp.String())
	}()

	status = http.StatusOK
}
