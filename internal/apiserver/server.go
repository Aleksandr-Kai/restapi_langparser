package apiserver

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"restapi_langparser/internal/apistructs"
	"restapi_langparser/internal/config"
	"restapi_langparser/internal/model"
	"restapi_langparser/internal/store"
)

type ctxKey int8

type server struct {
	router *gin.Engine
	store  store.IStore
	//finder *langfinder.LangFinder
}

func newServer(store store.IStore, config *config.Config) *server {
	s := &server{
		router: gin.New(),
		store:  store,
		//finder: langfinder.New(store, config),
	}
	s.configureRouter()
	return s
}

func (s *server) configureRouter() {
	s.router.Use(gin.Recovery())
	s.router.POST("/domains", s.handleAddDomains)
	s.router.GET("/domains", s.handleGetDomains)
	s.router.GET("/domains/:code", s.handleGetResult)
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

	res, err := s.store.GetRequest(c.Param("code"))
	if err != nil {
		resp.CreateErrorf(err.Error())
		return
	}
	resp.Results = &apistructs.APIResults{
		Domains: res,
	}
}

// handleGetDomains returns domains at the user's request
func (s *server) handleGetDomains(c *gin.Context) {
	var req apistructs.APIRequest
	resp := &apistructs.APIResponse{}
	status := http.StatusInternalServerError
	defer func() {
		c.String(status, resp.String())
	}()

	if err := c.BindJSON(&req); err != nil {
		resp.CreateErrorf("Internal error: %s", err.Error())
		status = http.StatusInternalServerError
		return
	}

	status = http.StatusOK

	resp.Results = &apistructs.APIResults{}
	domains, err := s.store.GetDomains(req.Hosts)
	if err != nil {
		status = http.StatusInternalServerError
		resp.CreateErrorf("Storage fail: %s", err.Error())
		return
	}

	resp.Results.RequestCode, err = s.store.CreateRequest(domains, req.Callback)
	if err != nil {
		status = http.StatusInternalServerError
		resp.CreateErrorf("Storage fail: %s", err.Error())
		return
	}

	if len(domains) == 0 {
		domains = model.CreateDomainsList(req.Hosts)
		if err = s.store.AddDomains(&domains); err != nil {
			status = http.StatusInternalServerError
			resp.CreateErrorf("Storage fail: %s", err.Error())
		}
		return
	}
	resp.Results.Domains = domains
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

	domains := model.CreateDomainsList(req.Hosts)
	if err := s.store.AddDomains(&domains); err != nil {
		status = http.StatusInternalServerError
		resp.CreateErrorf("Storage fail: %s", err.Error())
		return
	}
	status = http.StatusOK
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

	err := s.store.Proxy().Create(request.Proxy)
	if err != nil {
		resp.CreateErrorf(err.Error())
		return
	}

	status = http.StatusOK
}

func (s *server) handleGetProxyList(c *gin.Context) {
	resp := &apistructs.APIResponse{}
	status := http.StatusInternalServerError
	defer func() {
		c.String(status, resp.String())
	}()

	lst, err := s.store.Proxy().Read(0, 0)
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
