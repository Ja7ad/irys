package irys

import (
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
