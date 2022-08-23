package order

import (
	"context"
	"errors"
	"strconv"

	"github.com/sreway/gophermart/internal/entity"
	"github.com/sreway/gophermart/internal/usecase"
	"github.com/sreway/gophermart/pkg/utils"
)

type Order struct {
	order usecase.OrderRepo
}

func New(order usecase.OrderRepo) *Order {
	return &Order{
		order: order,
	}
}

func (oc *Order) Add(ctx context.Context, order *entity.Order) error {
	if err := order.Validate(); err != nil {
		return entity.ErrOrderIncorrectData
	}

	number, err := strconv.Atoi(order.Number)
	if err != nil {
		return entity.ErrOrderIncorrectNumber
	}

	if !utils.ValidLuhnNumber(number) {
		return entity.ErrOrderIncorrectNumber
	}

	err = oc.order.Add(ctx, order)
	if err != nil {
		if errors.Is(err, entity.ErrOrderAlreadyExist) {
			existOrder, getOrderErr := oc.order.Get(ctx, order.Number)
			switch {
			case getOrderErr != nil:
				return getOrderErr
			case existOrder.UserID != order.UserID:
				return entity.ErrOrderNumberTaken
			default:
				return entity.ErrOrderAlreadyExist
			}
		}
		return err
	}
	return nil
}

func (oc *Order) Get(ctx context.Context, userID uint) ([]*entity.Order, error) {
	orders, err := oc.order.GetAll(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(orders) == 0 {
		return nil, entity.ErrOrderEmptyData
	}
	return orders, err
}

func (oc *Order) UpdateStatus(ctx context.Context, number string, status entity.OrderStatus) error {
	return nil
}
