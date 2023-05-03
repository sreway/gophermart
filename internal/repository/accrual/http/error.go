package http

import (
	"fmt"
	"time"
)

type ErrRateLimited struct {
	RetryAfter time.Duration
}

func (eh *ErrRateLimited) Error() string {
	return fmt.Sprintf("Too Many Requests, try after %s", eh.RetryAfter)
}

func NewRateLimitError(retryAfter time.Duration) error {
	return &ErrRateLimited{
		RetryAfter: retryAfter,
	}
}
