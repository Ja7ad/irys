package irys

import "net/http"

type Option func(irys *Irys)

// WithCustomClient set custom http client for irys
func WithCustomClient(c *http.Client) Option {
	return func(irys *Irys) {
		irys.client = c
	}
}
