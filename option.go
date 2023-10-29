package irys

import (
	"net/http"
	"time"
)

type Option func(irys *IrysClient)

// WithCustomClient set custom http client for irys
func WithCustomClient(c *http.Client) Option {
	return func(irys *IrysClient) {
		irys.client.HTTPClient = c
	}
}

// WithCustomRetryMax maximum number of retries
func WithCustomRetryMax(retry int) Option {
	return func(irys *IrysClient) {
		irys.client.RetryMax = retry
	}
}

// WithCustomRetryWaitMin minimum time to wait
func WithCustomRetryWaitMin(waitMin time.Duration) Option {
	return func(irys *IrysClient) {
		irys.client.RetryWaitMin = waitMin
	}
}

// WithCustomRetryWaitMax maximum time to wait
func WithCustomRetryWaitMax(waitMax time.Duration) Option {
	return func(irys *IrysClient) {
		irys.client.RetryWaitMax = waitMax
	}
}
