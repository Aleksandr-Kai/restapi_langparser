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
}

type APIResults struct {
	Domains     []model.Domain `json:"Domains,omitempty"`
	Proxy       []model.Proxy  `json:"Proxy,omitempty"`
	RequestCode string         `json:"RequestCode,omitempty"`
}

type APIMessage struct {
	Text string
}

func CreateFromJSON(str string) (*APIResponse, error) {
	res := &APIResponse{}
	err := json.Unmarshal([]byte(str), res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *APIResponse) CreateErrorf(format string, arguments ...interface{}) {
	r.Error = &APIMessage{
		Text: fmt.Sprintf(format, arguments...),
	}
}

func (r *APIResponse) CreateMessagef(format string, arguments ...interface{}) {
	r.Message = &APIMessage{
		Text: fmt.Sprintf(format, arguments...),
	}
}

func (r *APIResponse) String() string {
	d, _ := json.Marshal(r)
	return string(d)
}
