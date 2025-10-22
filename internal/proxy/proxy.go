package proxy

import (
	"sync"
)

type ProxyPool struct {
	proxies []string
	index   int
	mu      sync.Mutex
}

func NewProxyPool(proxies []string) *ProxyPool {
	return &ProxyPool{
		proxies: proxies,
		index:   0,
	}
}

func (p *ProxyPool) GetNextProxy() string {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.proxies) == 0 {
		return ""
	}

	proxy := p.proxies[p.index]
	p.index = (p.index + 1) % len(p.proxies)
	return proxy
}

func (p *ProxyPool) GetProxyCount() int {
	return len(p.proxies)
}
