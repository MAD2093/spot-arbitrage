package proxy

import (
	"net/url"
	"sync"
)

type ProxyPool struct {
	proxies []*url.URL
	index   int
	mu      sync.RWMutex
}

func NewProxyPool(proxies []*url.URL) *ProxyPool {
	return &ProxyPool{
		proxies: proxies,
		index:   0,
	}
}

func (p *ProxyPool) GetNextProxy() (*url.URL, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.proxies) == 0 {
		return nil, false
	}

	proxy := p.proxies[p.index]
	p.index = (p.index + 1) % len(p.proxies)
	return proxy, true
}

func (p *ProxyPool) GetCurProxy() *url.URL {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.proxies[p.index]
}

func (p *ProxyPool) GetProxyCount() int {
	return len(p.proxies)
}
