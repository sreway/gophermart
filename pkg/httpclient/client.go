package httpclient

import (
	"github.com/go-resty/resty/v2"
)

type (
	Client struct {
		*resty.Client
	}
)

func New(opts ...Option) *Client {
	c := &Client{
		resty.New(),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}
