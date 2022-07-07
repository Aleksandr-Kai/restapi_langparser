package apistructs

import (
	"encoding/json"
	"fmt"
	"restapi_langparser/internal/model"
)

type APIResponse struct {
	Message *APIMessage `json:"Message,omitempty"`
	Error   *APIMessage `json:"Error,omitempty"`
	Results *APIResults `json:"Results,omitempty"`
	Status  int         `json:"-"`
}

type APIResults struct {
	Domains     []model.Domain `json:"Domains,omitempty"`
	Proxy       []model.Proxy  `json:"Proxy,omitempty"`
	RequestCode string         `json:"RequestCode,omitempty"`
}

type APIMessage string

func CreateFromJSON(str string) (*APIResponse, error) {
	res := &APIResponse{}
	err := json.Unmarshal([]byte(str), res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *APIResponse) CreateError(format string, args ...interface{}) {
	msg := APIMessage(fmt.Sprintf(format, args...))
	r.Error = &msg
}

func (r *APIResponse) CreateMessage(format string, args ...interface{}) {
	msg := APIMessage(fmt.Sprintf(format, args...))
	r.Message = &msg
}

func (r *APIResponse) String() string {
	d, _ := json.Marshal(r)
	return string(d)
}
