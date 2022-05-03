package apistructs

type APIRequest struct {
	Callback string   `json:"callback"`
	URLs     []string `json:"urls"`
}
