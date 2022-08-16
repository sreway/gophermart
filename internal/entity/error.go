package entity

import (
	"errors"
	"fmt"
	"time"
)

type (
	ErrApp struct {
		UserLogin string
		error     error
	}

	ErrUser struct {
		UserLogin string
		error     error
	}

	ErrOrder struct {
		UserLogin   string
		OrderNumber string
		error       error
	}

	ErrWithdraw struct {
		UserLogin   string
		OrderNumber string
		error       error
	}

	ErrHTTPClient struct {
		StatusCode int
		error      error
	}

	ErrRateLimited struct {
		RetryAfter time.Duration
	}
)

func NewUserError(userLogin string, err error) error {
	return &ErrUser{
		UserLogin: userLogin,
		error:     err,
	}
}

func NewAppError(userLogin string, err error) error {
	return &ErrApp{
		UserLogin: userLogin,
		error:     err,
	}
}

func NewOrderError(userLogin, orderNumber string, err error) error {
	return &ErrOrder{
		UserLogin:   userLogin,
		OrderNumber: orderNumber,
		error:       err,
	}
}

func NewErrWithdraw(userLogin, orderNumber string, err error) error {
	return &ErrWithdraw{
		UserLogin:   userLogin,
		OrderNumber: orderNumber,
		error:       err,
	}
}

func NewErrHTTPClient(statusCode int, err error) error {
	return &ErrHTTPClient{
		StatusCode: statusCode,
		error:      err,
	}
}

func NewRateLimitError(retryAfter time.Duration) error {
	return &ErrRateLimited{
		RetryAfter: retryAfter,
	}
}

func (eu *ErrUser) Error() string {
	return fmt.Sprintf("User_Error[%s]: %s", eu.UserLogin, eu.error)
}

func (eu *ErrUser) Is(err error) bool {
	return errors.Is(eu.error, err)
}

func (ea *ErrApp) Error() string {
	return fmt.Sprintf("App_Error[%s]: %s", ea.UserLogin, ea.error)
}

func (ea *ErrApp) Is(err error) bool {
	return errors.Is(ea.error, err)
}

func (eo *ErrOrder) Error() string {
	return fmt.Sprintf("Order_Error[%s][%s]: %s", eo.UserLogin, eo.OrderNumber, eo.error)
}

func (eo *ErrOrder) Is(err error) bool {
	return errors.Is(eo.error, err)
}

func (ew *ErrWithdraw) Error() string {
	return fmt.Sprintf("Withdraw_Error[%s][%s]: %s", ew.UserLogin, ew.OrderNumber, ew.error)
}

func (ew *ErrWithdraw) Is(err error) bool {
	return errors.Is(ew.error, err)
}

func (eh *ErrHTTPClient) Error() string {
	return fmt.Sprintf("HTTPClient_Error[%d]: %s", eh.StatusCode, eh.error)
}

func (eh *ErrRateLimited) Error() string {
	return fmt.Sprintf("RateLimited_Error: Too Many Requests, try after %s", eh.RetryAfter)
}

var (
	ErrUserIncorrectData     = fmt.Errorf("incorrect user data")
	ErrUserAlreadyExist      = fmt.Errorf("user already exists")
	ErrUserUnauthorized      = fmt.Errorf("user unauthorized")
	ErrUserNotFound          = fmt.Errorf("user not found")
	ErrUserIncorrectPassword = fmt.Errorf("incorrect user password")
	ErrOrderIncorrectData    = fmt.Errorf("incorrect order data")
	ErrOrderAlreadyExist     = fmt.Errorf("order number already exist")
	ErrOrderNumberTaken      = fmt.Errorf("order number already taken other user")
	ErrOrderIncorrectNumber  = fmt.Errorf("incorrect order number")
	ErrOrderEmptyData        = fmt.Errorf("empty data")
	ErrOrderNotFound         = fmt.Errorf("order not found")
	ErrBalanceNotEnough      = fmt.Errorf("not enough funds on the balance")
	ErrWithdrawEmptyData     = fmt.Errorf("empty data")
)
