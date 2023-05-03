package domain

import (
	"errors"
	"fmt"
)

var (
	ErrIncorrectData        = errors.New("incorrect data")
	ErrAlreadyExist         = errors.New("already exists")
	ErrUserUnauthorized     = errors.New("user unauthorized")
	ErrNotFound             = errors.New("not found")
	ErrOrderNumberInvalid   = errors.New("invalid order number")
	ErrOrderStatusInvalid   = errors.New("invalid order status")
	ErrAccrualStatusInvalid = errors.New("invalid accrual status")
	ErrRateLimit            = errors.New("error rate limit")
	ErrBalanceNotEnough     = errors.New("not enough funds on the balance")
	ErrOrderNumberTaken     = fmt.Errorf("order number already taken other user")
)

type (
	ErrUser struct {
		id    string
		error error
	}

	ErrOrder struct {
		id    string
		error error
	}
	ErrAccrual struct {
		number int
		error  error
	}
	ErrBalance struct {
		userID string
		error  error
	}

	ErrWithdraw struct {
		id    string
		error error
	}
)

func (eu *ErrUser) Error() string {
	return eu.error.Error()
}

func (eu *ErrUser) UserID() string {
	return eu.id
}

func (eu *ErrUser) Is(err error) bool {
	return errors.Is(eu.error, err)
}

func (eo *ErrOrder) Error() string {
	return eo.error.Error()
}

func (eo *ErrOrder) OrderID() string {
	return eo.id
}

func (eo *ErrOrder) Is(err error) bool {
	return errors.Is(eo.error, err)
}

func (ea *ErrAccrual) Error() string {
	return ea.error.Error()
}

func (ea *ErrAccrual) Is(err error) bool {
	return errors.Is(ea.error, err)
}

func (ea *ErrAccrual) OrderNumber() string {
	return fmt.Sprintf("%d", ea.number)
}

func (eb *ErrBalance) Error() string {
	return eb.error.Error()
}

func (eb *ErrBalance) Is(err error) bool {
	return errors.Is(eb.error, err)
}

func (eb *ErrBalance) UserID() string {
	return eb.userID
}

func (ew *ErrWithdraw) Error() string {
	return ew.error.Error()
}

func (ew *ErrWithdraw) Is(err error) bool {
	return errors.Is(ew.error, err)
}

func (ew *ErrWithdraw) ID() string {
	return ew.id
}

func NewUserError(id string, err error) error {
	return &ErrUser{
		id:    id,
		error: err,
	}
}

func NewOrderError(id string, err error) error {
	return &ErrUser{
		id:    id,
		error: err,
	}
}

func NewAccrualError(number int, err error) error {
	return &ErrAccrual{
		number: number,
		error:  err,
	}
}

func NewBalanceError(userID string, err error) error {
	return &ErrBalance{
		userID: userID,
		error:  err,
	}
}

func NewWithdrawError(id string, err error) error {
	return &ErrWithdraw{
		id:    id,
		error: err,
	}
}
