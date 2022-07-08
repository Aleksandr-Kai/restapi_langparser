package apiserver

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"restapi_langparser/internal/apistructs"
	"restapi_langparser/internal/config"
	"restapi_langparser/internal/model"
	"restapi_langparser/internal/store"
	"strconv"
	"strings"
)

type ctxKey int8

type server struct {
	router *gin.Engine
	store  store.IStore
	config config.Config
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

func newResp(c *gin.Context) (*apistructs.APIResponse, func()) {
	resp := &apistructs.APIResponse{
		Status: http.StatusOK,
	}
	return resp, func() {
		c.String(resp.Status, resp.String())
	}
}

func (s *server) getDomainByID(sid string) (*apistructs.APIResults, error) {
	id, err := strconv.Atoi(sid)
	if err != nil {
		return nil, err
	}

	domain, err := s.store.Domain().FindByID(id)
	if err != nil {
		return nil, err
	}
	return &apistructs.APIResults{
		Domains: []model.Domain{
			*domain,
		},
	}, nil
}

func (s *server) requestDomains(hostsString string, callback *string) (*apistructs.APIResults, error) {
	hosts := strings.Split(hostsString, ",")

	domains, err := s.store.GetDomains(hosts)
	if err != nil {
		return nil, err
	}

	res := &apistructs.APIResults{}

	if len(domains) == 0 {
		domains = model.CreateDomainsList(hosts)
		if err = s.store.AddDomains(&domains); err != nil {
			return nil, err
		}
	} else {
		res.Domains = domains
	}

	requestCode, err := s.store.CreateRequest(domains, callback)
	if err != nil {
		return nil, err
	}

	res.RequestCode = requestCode

	return res, nil
}

func (s *server) configureRouter() {
	s.router.Use(gin.Recovery())
	s.router.POST("/domains", s.handleAddDomains)
	s.router.GET("/domains", s.domainsHandler)
	s.router.PUT("/domains/:id", s.handleUpdateDomain)
	s.router.DELETE("/domains/:id", s.handleDeleteDomain)

	s.router.POST("/proxy", s.handleAddProxy)
	s.router.GET("/proxy", s.handleGetProxyList)
	s.router.PUT("/proxy/:id", s.handleUpdateProxy)
	s.router.DELETE("/proxy/:id", s.handleDeleteProxy)

	s.router.POST("/echo", func(c *gin.Context) {
		resp, writeResp := newResp(c)
		defer writeResp()

		d, _ := io.ReadAll(c.Request.Body)
		c.Request.Body.Close()
		fmt.Println(string(d))
		resp.CreateMessage("Echo request accepted")
	})
}

func (s *server) Start(addr string) error {
	return s.router.Run(addr)
}

/*
	Handlers
*/

func (s *server) domainsHandler(c *gin.Context) {
	resp, writeResp := newResp(c)
	defer writeResp()

	var err error

	if id := c.Query("id"); id != "" { // get domain info by ID
		resp.Results, err = s.getDomainByID(id)
		if err != nil {
			resp.Status = http.StatusInternalServerError
			resp.CreateError(err.Error())
			return
		}
	} else if hosts := c.Query("hosts"); len(hosts) > 0 { // get domains info by host
		var cb *string
		callback := c.Query("callback")
		if callback != "" {
			cb = &callback
		}

		resp.Results, err = s.requestDomains(hosts, cb)
		if err != nil {
			resp.Status = http.StatusInternalServerError
			resp.CreateError(err.Error())
			return
		}
	} else if code := c.Query("code"); code != "" { // get results by request code
		res, err := s.store.GetRequest(code)
		if err != nil {
			resp.Status = http.StatusInternalServerError
			resp.CreateError(err.Error())
			return
		}
		resp.Results = &apistructs.APIResults{
			Domains: res,
		}
	} else { // list all domains
		limit, err := strconv.Atoi(c.DefaultQuery("pagesize", "10"))
		if err != nil {
			resp.Status = http.StatusInternalServerError
			resp.CreateError(err.Error())
			return
		}
		page, err := strconv.Atoi(c.DefaultQuery("page", "0"))
		if err != nil {
			resp.Status = http.StatusInternalServerError
			resp.CreateError(err.Error())
			return
		}

		res, err := s.store.Domain().Read(limit, page*limit)
		if err != nil {
			resp.Status = http.StatusInternalServerError
			resp.CreateError(err.Error())
			return
		}
		resp.Results = &apistructs.APIResults{
			Domains: res,
		}
		resp.CreateMessage("Page size: %d, Page: %d", limit, page)
	}
}

func (s *server) handleUpdateDomain(c *gin.Context) {
	var domain model.Domain
	resp := &apistructs.APIResponse{}
	status := http.StatusInternalServerError
	defer func() {
		c.String(status, resp.String())
	}()

	if err := c.BindJSON(&domain); err != nil {
		resp.CreateError(err.Error())
		status = http.StatusInternalServerError
		return
	}

	status = http.StatusOK

	sID := c.Param("id")
	id, _ := strconv.Atoi(sID)
	domain.ID = uint(id)

	if err := s.store.Domain().Update(domain); err != nil {
		resp.Status = http.StatusInternalServerError
		resp.CreateError(err.Error())
		return
	}
	resp.CreateMessage("updated")
}

func (s *server) handleDeleteDomain(c *gin.Context) {
	resp := &apistructs.APIResponse{}
	status := http.StatusInternalServerError
	defer func() {
		c.String(status, resp.String())
	}()

	status = http.StatusOK
	resp.CreateMessage("request %v", c.Param("id"))
}

func (s *server) handleAddDomains(c *gin.Context) {
	var req apistructs.APIRequest
	resp := &apistructs.APIResponse{}
	status := http.StatusInternalServerError
	defer func() {
		c.String(status, resp.String())
	}()

	body, err := io.ReadAll(c.Request.Body)
	c.Request.Body.Close()
	if err != nil {
		resp.CreateError(err.Error())
		status = http.StatusInternalServerError
		return
	}

	if err = json.Unmarshal(body, &req); err != nil {
		resp.CreateError(err.Error())
		status = http.StatusInternalServerError
		return
	}

	if len(req.Hosts) == 0 {
		logrus.Errorf("Invalid data format in the request %s", string(body))
		resp.CreateError("Invalid data format in the request %s", string(body))
		status = http.StatusInternalServerError
		return
	}

	domains := model.CreateDomainsList(req.Hosts)
	if err = s.store.AddDomains(&domains); err != nil {
		status = http.StatusInternalServerError
		resp.CreateError("Storage fail: %s", err.Error())
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
		resp.CreateError(err.Error())
		status = http.StatusInternalServerError
		return
	}

	err := s.store.Proxy().Create(request.Proxy)
	if err != nil {
		resp.CreateError(err.Error())
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
		resp.CreateError(err.Error())
		return
	}
	resp.Results = &apistructs.APIResults{
		Proxy: lst,
	}
	resp.CreateMessage("proxy %v", lst)

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
