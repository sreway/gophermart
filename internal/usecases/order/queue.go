package order

import (
	"context"
	"errors"
	"time"

	"golang.org/x/exp/slog"

	"github.com/sreway/gophermart/internal/config"
	"github.com/sreway/gophermart/internal/domain"
	"github.com/sreway/gophermart/internal/repository/accrual/http"
)

const (
	procAccrual action = "accrual processing"
)

type (
	action string
	task   struct {
		name  action
		order domain.Order
	}
)

func NewTask(name action, order domain.Order) *task {
	return &task{
		name:  name,
		order: order,
	}
}

func (uc *useCase) procTaskAccrual(ctx context.Context, order domain.Order) error {
	accrual, err := uc.accrual.Get(ctx, *order.Number())
	if err != nil {
		var accrualError *http.ErrRateLimited
		switch {
		case errors.As(err, &accrualError):
			time.Sleep(accrualError.RetryAfter)
			uc.tasks <- NewTask(procAccrual, order)
			return err
		default:
			uc.logger.Error("failed get accrual status", slog.Any("err", err),
				slog.Int("number", order.Number().Value()))
			uc.tasks <- NewTask(procAccrual, order)
			return err
		}
	}
	var procOrder *domain.Order

	switch accrual.Status() {
	case domain.AccrualProcessed:
		procOrder, err = uc.storage.UpdateOrderStatus(ctx, *order.Number(), domain.OrderProcessed, accrual.Value())
		if err != nil {
			uc.logger.Error("failed update order status", slog.Any("err", err),
				slog.Int("number", order.Number().Value()))
			return err
		}

		err = uc.storage.RefillBalance(ctx, procOrder.UserID(), accrual.Value())
		if err != nil {
			uc.logger.Error("failed update user balance", slog.Any("err", err),
				slog.String("user_id", procOrder.UserID().String()))
			return err
		}

		uc.logger.Info("success update user balance", slog.String("user_id", procOrder.UserID().String()))

	case domain.AccrualInvalid:
		_, err = uc.storage.UpdateOrderStatus(ctx, *order.Number(), domain.OrderInvalid, accrual.Value())
		if err != nil {
			uc.logger.Error("failed update order status", slog.Any("err", err),
				slog.Int("number", order.Number().Value()))
			return err
		}
	case domain.AccrualProcessing:
		procOrder, err := uc.storage.UpdateOrderStatus(ctx, *order.Number(), domain.OrderProcessing, accrual.Value())
		if err != nil {
			uc.logger.Error("failed update order status", slog.Any("err", err))
			return err
		}
		uc.tasks <- NewTask(procAccrual, *procOrder)
	default:
		uc.tasks <- NewTask(procAccrual, order)
	}

	return nil
}

func (uc *useCase) ProcNewOrder(ctx context.Context, config *config.Orders) error {
	tick := time.NewTicker(config.TaskInterval)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			if len(uc.tasks) == 0 {
				continue
			}

			actions := make(map[action][]domain.Order, len(uc.tasks))

			for len(uc.tasks) != 0 {
				t := <-uc.tasks
				actions[t.name] = append(actions[t.name], t.order)
			}

			for k, v := range actions {
				switch k {
				case procAccrual:
					for _, item := range v {
						order := item
						go func() {
							err := uc.procTaskAccrual(ctx, order)
							if err != nil {
								uc.logger.Error("failed processing order in accrual system",
									slog.Any("err", err))
								return
							}
							uc.logger.Info("success processing order in accrual system",
								slog.Int("number", order.Number().Value()))
						}()
					}

				default:
					uc.logger.Warn("unknown task action", slog.Any("action", k),
						slog.String("func", "ProcNewOrder"))
				}
			}
		case <-ctx.Done():
			close(uc.tasks)
			uc.logger.Info("stop processed new orders queue")
			return nil
		}
	}
}
