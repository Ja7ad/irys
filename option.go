package irys

import (
	"golang.org/x/net/proxy"
	"net/http"
	"time"

	"github.com/Ja7ad/irys/utils/logger"
)

type Option func(irys *Client)

// WithCustomClient set custom http client for irys
func WithCustomClient(c *http.Client) Option {
	return func(irys *Client) {
		irys.client.HTTPClient = c
	}
}

// WithCustomRetryMax maximum number of retries
func WithCustomRetryMax(retry int) Option {
	return func(irys *Client) {
		irys.client.RetryMax = retry
	}
}

// WithCustomRetryWaitMin minimum time to wait
func WithCustomRetryWaitMin(waitMin time.Duration) Option {
	return func(irys *Client) {
		irys.client.RetryWaitMin = waitMin
	}
}

// WithCustomRetryWaitMax maximum time to wait
func WithCustomRetryWaitMax(waitMax time.Duration) Option {
	return func(irys *Client) {
		irys.client.RetryWaitMax = waitMax
	}
}

// WithCustomLogging create custom logging
func WithCustomLogging(logging logger.Logger) Option {
	return func(irys *Client) {
		irys.logging = logging
	}
}

// WithHttpProxy add http proxy for client
//
// Example:
//
//	c, err := irys.New(irys.DefaultNode1, matic, true, irys.WithHttpProxy("http://foobar.com:8080"))
func WithHttpProxy(uri string) Option {
	return func(irys *Client) {
		irys.proxy.proxyType = _http
		irys.proxy.uri = uri
	}
}

// WithSocks5Proxy add socks5 proxy for client
//
// Example:
//
//	c, err := irys.New(irys.DefaultNode1, matic, true, irys.WithSocks5Proxy("ip:port", "foo", "bar"))
func WithSocks5Proxy(ip, username, password string) Option {
	return func(irys *Client) {
		irys.proxy.proxyType = _socks5
		irys.proxy.uri = ip
		if len(username) != 0 && len(password) != 0 {
			irys.proxy.auth = proxy.Auth{
				User:     username,
				Password: password,
			}
		}
	}
}
