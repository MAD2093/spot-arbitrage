package req

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"golang.org/x/time/rate"
	"mad-scanner.com/scriner/internal/proxy"
)

type ExchangeClient struct {
	name         string
	limiter      *rate.Limiter
	proxyPool    *proxy.ProxyPool
	mu           sync.Mutex
	failCount    int
	currentProxy string
}

type ClientConfig struct {
	Name      string
	RPS       int // requests per second
	Burst     int
	ProxyPool *proxy.ProxyPool
}

func NewExchangeClient(config ClientConfig) *ExchangeClient {
	return &ExchangeClient{
		name:      config.Name,
		limiter:   rate.NewLimiter(rate.Limit(config.RPS), config.Burst),
		proxyPool: config.ProxyPool,
		failCount: 0,
	}
}

func (c *ExchangeClient) MakeRequest(req *http.Request) ([]byte, error) {
	// Wait for rate limit
	ctx := context.Background()
	if err := c.limiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("%s rate limit error: %w", c.name, err)
	}

	var body []byte
	var err error

	// Try with current proxy or get new one
	for attempt := 0; attempt < 3; attempt++ {
		if c.currentProxy == "" || c.failCount >= 2 {
			c.rotateProxy()
		}

		body, err = c.doRequestWithProxy(req)
		if err == nil {
			c.failCount = 0 // reset fail count on success
			return body, nil
		}

		c.failCount++
		fmt.Printf("%s request failed (attempt %d): %v\n", c.name, attempt+1, err)
	}

	return nil, fmt.Errorf("%s all attempts failed: %w", c.name, err)
}

func (c *ExchangeClient) rotateProxy() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.currentProxy = c.proxyPool.GetNextProxy()
	c.failCount = 0
	fmt.Printf("%s rotated to proxy: %s\n", c.name, c.currentProxy)
}

func (c *ExchangeClient) doRequestWithProxy(req *http.Request) ([]byte, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Set proxy if available
	if c.currentProxy != "" {
		proxyURL, err := url.Parse(c.currentProxy)
		if err != nil {
			return nil, fmt.Errorf("invalid proxy URL: %w", err)
		}
		client.Transport = &http.Transport{
			Proxy:               http.ProxyURL(proxyURL),
			TLSHandshakeTimeout: 10 * time.Second,
			IdleConnTimeout:     30 * time.Second,
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	return io.ReadAll(resp.Body)
}

// requestBytes attempts to perform the request using a configured ExchangeClient
// (which may add rate-limiting and proxy handling). If no client is available
// it falls back to the default http client.
func requestBytes(req *http.Request, exchange string) ([]byte, error) {
	client := GetExchangeClient(exchange)
	if client != nil {
		return client.MakeRequest(req)
	}

	// fallback
	c := &http.Client{Timeout: 10 * time.Second}
	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	return io.ReadAll(resp.Body)
}
