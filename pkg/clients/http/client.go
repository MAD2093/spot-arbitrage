package http

import (
	"errors"
	"io"
	"net/http"
	"net/url"

	rdmutex "github.com/Q1rD/rdmutex/rdmutex"

	"github.com/rs/zerolog"
	"mad-scanner.com/scriner/pkg/clients/http/proxy"
	"mad-scanner.com/scriner/pkg/log"
)

type HTTPClient struct {
	PP       *proxy.ProxyPool
	curProxy *url.URL

	log *zerolog.Logger

	rdmu *rdmutex.RDMutex
}

func NewHTTPClient(pp *proxy.ProxyPool) *HTTPClient {
	log := log.GetLogger()

	return &HTTPClient{
		PP:       pp,
		curProxy: pp.GetCurProxy(),
		log:      log,
		rdmu:     rdmutex.NewRDMutex(),
	}
}

func (c *HTTPClient) MakeRequest(req *http.Request) ([]byte, error) {
	c.rdmu.RLock()
	proxy := c.curProxy

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxy),
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	need := c.needToRotateProxy(resp)
	if need {
		c.rdmu.RUnlock()
		for need {
			if c.rdmu.Lock() {
				// писатель
				// получаем новый прокси
				if c.rotateProxy() {
					return nil, errors.New("fail while trying to rotate proxy")
				}

				// переназначаем прокси у клиента
				client.Transport = &http.Transport{
					Proxy: http.ProxyURL(c.curProxy),
				}
				// делаем запрос с новым прокси
				resp, err = client.Do(req)
				if err != nil {
					return nil, err
				}
				need = c.needToRotateProxy(resp)
				c.rdmu.Unlock()
			} else {
				// читатель
				c.rdmu.Wait()

				// переназначаем прокси
				client.Transport = &http.Transport{
					Proxy: http.ProxyURL(c.curProxy),
				}
				// делаем запрос с новым прокси
				resp, err = client.Do(req)
				if err != nil {
					return nil, err
				}
				need = c.needToRotateProxy(resp)

				// освобождаем читателя
				c.rdmu.RUnlock()
			}
		}
		return io.ReadAll(resp.Body)
	} else {
		c.rdmu.RUnlock()
		return io.ReadAll(resp.Body)
	}
}

func (c *HTTPClient) needToRotateProxy(resp *http.Response) bool {
	// client errors
	if 400 <= resp.StatusCode && resp.StatusCode < 500 {
		return true
	}
	//  server errors
	if 500 <= resp.StatusCode && resp.StatusCode < 600 {
		return false
	}
	return false
}

func (c *HTTPClient) rotateProxy() bool {
	proxy, ok := c.PP.GetNextProxy()
	if !ok {
		c.log.Warn().Msg("fail while trying to rotate proxy")
		return false
	}
	c.curProxy = proxy
	return true
}
