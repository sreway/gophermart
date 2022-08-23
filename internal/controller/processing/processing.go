package processing

import (
	"context"
	"errors"
	"time"

	"github.com/sreway/gophermart/internal/entity"
	"github.com/sreway/gophermart/internal/usecase"
	"github.com/sreway/gophermart/pkg/logger"
)

type Orders struct {
	queue   usecase.Queue
	accrual usecase.Accrual
}

func NewOrders(q usecase.Queue, a usecase.Accrual) (*Orders, error) {
	return &Orders{
		queue:   q,
		accrual: a,
	}, nil
}

func (po *Orders) SendToQueue(ctx context.Context, data <-chan string) {
	for {
		select {
		case <-ctx.Done():
			return
		case orderNumber := <-data:
			err := po.queue.Add(ctx, orderNumber)
			if err != nil {
				logger.Error(err)
			}
		}
	}
}

func (po *Orders) CheckInAccrual(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := po.queue.Read(ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					return
				}
				logger.Error(err)
				continue
			}

			if msg == nil {
				continue
			}

			a, err := po.accrual.Get(ctx, string(msg.Value))
			if err != nil {
				var rateLimitedError *entity.ErrRateLimited
				var httpClientError *entity.ErrHTTPClient

				switch {
				case errors.As(err, &rateLimitedError):
					logger.Error(rateLimitedError)
					time.Sleep(rateLimitedError.RetryAfter)
					continue
				case errors.As(err, &httpClientError):
					logger.Error(httpClientError)
					continue
				}
			}

			err = po.accrual.UpdateOrderStatus(ctx, a)
			if err != nil {
				logger.Error(err)
				continue
			}

			switch entity.OrderStatus(a.Status) {
			case entity.OrderStatusInvalid, entity.OrdertStatusProcessed:
				logger.Infof("Success processed order %s", string(msg.Value))
			default:
				err = po.queue.Commit(ctx, msg)
				if err != nil {
					logger.Error(err)
					continue
				}
			}
		}
	}
}
