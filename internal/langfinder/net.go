package langfinder

import (
	"net/http"
	"net/url"
	"restapi_langparser/internal/model"
	"time"
)

func createClient(proxy *model.Proxy, timeout time.Duration) *http.Client {
	switch proxy.Type {
	case model.HTTPS:
		proxyURL := &url.URL{
			Scheme: "https://",
			Host:   proxy.IP + ":" + proxy.Port,
		}
		if proxy.Login != "" {
			proxyURL.Host = proxy.Login + ":" + proxy.Password + "@" + proxyURL.Host
		}
		return &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			},
		}
	case model.Socks5:
	case model.Socks4:
	}
	return &http.Client{
		Timeout: timeout,
	}
}
