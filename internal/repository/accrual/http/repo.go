package http

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/go-resty/resty/v2"
	"golang.org/x/exp/slog"

	"github.com/sreway/gophermart/internal/config"
	"github.com/sreway/gophermart/internal/domain"
)

type (
	repo struct {
		client *resty.Client
		logger *slog.Logger
	}

	accrualResponse struct {
		OrderNumber string  `json:"order"`
		Status      string  `json:"status"`
		Value       float64 `json:"accrual"`
	}
)

func (r *repo) Get(ctx context.Context, number domain.OrderNumber) (*domain.Accrual, error) {
	endpoint := fmt.Sprintf("/api/orders/%d", number.Value())
	data := new(accrualResponse)

	response, err := r.client.R().SetContext(ctx).SetResult(data).Get(endpoint)
	if err != nil {
		r.logger.Error("failed send request for check accrual", slog.Any("err", err))
		return nil, domain.NewAccrualError(number.Value(), err)
	}

	switch response.StatusCode() {
	case http.StatusOK:
		accrual, err := domain.NewAccrual(data.OrderNumber, data.Status, data.Value)
		if err != nil {
			r.logger.Error("failed get accrual", slog.Any("err", err))
			return nil, domain.NewAccrualError(number.Value(), err)
		}
		return accrual, nil
	case http.StatusTooManyRequests:
		return nil, domain.NewAccrualError(number.Value(), domain.ErrRateLimit)
	case http.StatusNoContent:
		return nil, domain.NewAccrualError(number.Value(), domain.ErrNotFound)
	default:
		return nil, domain.NewAccrualError(number.Value(), fmt.Errorf(string(response.Body())))
	}
}

func New(cfg *config.Accrual) *repo {
	log := slog.New(slog.NewJSONHandler(os.Stdout).
		WithAttrs([]slog.Attr{slog.String("repository", "accrual")}))
	client := resty.New().SetBaseURL(cfg.Address)

	return &repo{
		client: client,
		logger: log,
	}
}
