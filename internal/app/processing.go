package app

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/sreway/gophermart/config"
	"github.com/sreway/gophermart/internal/controller/listener"
	"github.com/sreway/gophermart/internal/entity"
	"github.com/sreway/gophermart/internal/usecase/accrual"
	"github.com/sreway/gophermart/internal/usecase/repo"
	"github.com/sreway/gophermart/pkg/logger"
	"github.com/sreway/gophermart/pkg/postgres"
)

func orderListener(ctx context.Context, wg *sync.WaitGroup, pg *postgres.Postgres, cfg *config.Config, data chan<- string) {
	defer func() {
		logger.Info("Stop orderListener")
		wg.Done()
	}()
	pgl, err := listener.NewPgListner(ctx, pg, cfg.Postgres.ListenChannel)
	if err != nil {
		logger.Fatal(err)
	}
	defer pgl.Release()

	pgl.Listen(ctx, data)
}

func storeOrderQueue(ctx context.Context, wg *sync.WaitGroup, oq *repo.QueueRepo, data <-chan string) {
	defer func() {
		logger.Info("Stop storeOrderInKafka")
		wg.Done()
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case orderNumber := <-data:
			err := oq.Store(orderNumber)
			if err != nil {
				logger.Error(err)
			}
		}
	}
}

func processingOrder(ctx context.Context, wg *sync.WaitGroup, ac *accrual.Accrual, oq *repo.QueueRepo) {
	defer func() {
		logger.Info("Stop processingOrder")
		wg.Done()
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := oq.Read(ctx)
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

			a, err := ac.Get(ctx, string(msg.Value))
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

			err = ac.UpdateOrderStatus(ctx, a)
			if err != nil {
				logger.Error(err)
				continue
			}

			switch entity.OrderStatus(a.Status) {
			case entity.OrderStatusInvalid, entity.OrdertStatusProcessed:
				logger.Infof("Success processed order %s", string(msg.Value))
			default:
				err = oq.Commit(ctx, msg)
				if err != nil {
					logger.Error(err)
					continue
				}
			}
		}
	}
}
