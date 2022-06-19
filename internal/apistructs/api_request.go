package apistructs

import "restapi_langparser/internal/model"

type APIRequest struct {
	Callback *string       `json:"callback,omitempty"`
	Hosts    []string      `json:"urls,omitempty"`
	Proxy    []model.Proxy `json:"proxy,omitempty"`
}
