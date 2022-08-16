package repo

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/sreway/gophermart/internal/entity"
)

type AccrualRepo struct {
	*resty.Client
}

func NewAccrualRepo(c *resty.Client) *AccrualRepo {
	return &AccrualRepo{c}
}

func (ar *AccrualRepo) Get(ctx context.Context, number string) (*entity.Accrual, error) {
	endpoint := fmt.Sprintf("/api/orders/%s", number)
	a := entity.Accrual{}
	r, err := ar.R().SetContext(ctx).SetResult(&a).Get(endpoint)
	if err != nil {
		return nil, entity.NewErrHTTPClient(r.StatusCode(), err)
	}

	switch r.StatusCode() {
	case http.StatusOK:
		return &a, err
	case http.StatusTooManyRequests:
		retry, rErr := strconv.Atoi(r.Header().Get("Retry-After"))
		if rErr != nil {
			return nil, entity.NewErrHTTPClient(r.StatusCode(), rErr)
		}
		return nil, entity.NewRateLimitError(time.Second * time.Duration(retry))
	case http.StatusNoContent:
		return nil, entity.NewErrHTTPClient(r.StatusCode(), entity.ErrOrderNotFound)
	default:
		return nil, entity.NewErrHTTPClient(r.StatusCode(), fmt.Errorf(string(r.Body())))
	}
}
