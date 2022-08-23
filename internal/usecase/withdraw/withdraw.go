package withdraw

import (
	"context"
	"encoding/base64"
	"strconv"

	"github.com/sreway/gophermart/internal/entity"
	"github.com/sreway/gophermart/internal/usecase"
)

const DefaultMaxResult = 10

type Withdraw struct {
	withdraw usecase.WithdrawRepo
}

func New(withdraw usecase.WithdrawRepo) *Withdraw {
	return &Withdraw{
		withdraw: withdraw,
	}
}

func (wc *Withdraw) Add(ctx context.Context, withdraw *entity.Withdraw) error {
	err := wc.withdraw.Add(ctx, withdraw)
	if err != nil {
		return err
	}

	return nil
}

func (wc *Withdraw) Get(ctx context.Context, userID uint) ([]*entity.WithdrawOrder, error) {
	withdrawals, err := wc.withdraw.Get(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(withdrawals) == 0 {
		return nil, entity.ErrWithdrawEmptyData
	}

	return withdrawals, err
}

func (wc *Withdraw) GetPagination(ctx context.Context, userID uint, nextPageToken,
	pageSize string,
) (*entity.Withdrawals, error) {
	var startAt, maxResults int

	if len(nextPageToken) != 0 {
		data, err := base64.StdEncoding.DecodeString(nextPageToken)
		if err != nil {
			return nil, err
		}

		startAt, err = strconv.Atoi(string(data))
		if err != nil {
			return nil, err
		}
	}

	if len(pageSize) == 0 {
		maxResults = DefaultMaxResult
	} else {
		var err error
		maxResults, err = strconv.Atoi(pageSize)
		if err != nil {
			return nil, err
		}
	}

	withdrawals, err := wc.withdraw.GetPagination(ctx, userID, uint(startAt), uint(maxResults))
	if err != nil {
		return nil, err
	}

	if len(withdrawals.Items) == 0 {
		return nil, entity.ErrWithdrawEmptyData
	}

	lastID := withdrawals.Items[len(withdrawals.Items)-1].ID

	if !(len(withdrawals.Items) < int(withdrawals.PageSize)) {
		withdrawals.NextPageToken = base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(int(lastID))))
	}

	return withdrawals, nil
}
