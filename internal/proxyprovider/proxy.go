package proxyprovider

type Proxy struct {
	threadLimit chan interface{}
}

func NewProxy(maxThreadCount uint) *Proxy {
	return &Proxy{
		threadLimit: make(chan interface{}, maxThreadCount-1),
	}
}
