package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type (
	OrderStatus string
	OrderNumber struct {
		value int
	}
	Order struct {
		id         uuid.UUID
		number     OrderNumber
		userID     uuid.UUID
		status     OrderStatus
		accrual    float64
		uploadedAt time.Time
	}
)

const (
	OrderNew        OrderStatus = "NEW"
	OrderProcessing OrderStatus = "PROCESSING"
	OrderInvalid    OrderStatus = "INVALID"
	OrderProcessed  OrderStatus = "PROCESSED"
)

func (o *Order) ID() uuid.UUID {
	return o.id
}

func (o *Order) Number() *OrderNumber {
	return &o.number
}

func (o *Order) UserID() uuid.UUID {
	return o.userID
}

func (o *Order) Status() OrderStatus {
	return o.status
}

func (o *Order) Accrual() float64 {
	return o.accrual
}

func (o *Order) UploadedAt() time.Time {
	return o.uploadedAt
}

func (on *OrderNumber) Value() int {
	return on.value
}

func (os OrderStatus) String() string {
	return string(os)
}

func (os OrderStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(os.String())
}

func NewOrderNumber(number int) (*OrderNumber, error) {
	checksum := func(number int) int {
		var sum int

		for i := 0; number > 0; i++ {
			cur := number % 10

			if i%2 == 0 {
				cur *= 2
				if cur > 9 {
					cur = cur%10 + cur/10
				}
			}

			sum += cur
			number /= 10
		}
		return sum % 10
	}

	if !((number%10+checksum(number/10))%10 == 0) {
		return nil, ErrOrderNumberInvalid
	}

	return &OrderNumber{
		value: number,
	}, nil
}

func NewOrderStatus(value string) (*OrderStatus, error) {
	status := OrderStatus(value)
	switch status {
	case OrderNew, OrderInvalid, OrderProcessing, OrderProcessed:
		return &status, nil
	default:
		return nil, ErrOrderStatusInvalid
	}
}

func NewOrder(userID uuid.UUID, number OrderNumber, status OrderStatus) *Order {
	id := uuid.New()
	return &Order{
		id:     id,
		userID: userID,
		number: number,
		status: status,
	}
}

func CreateOrder(id, userID uuid.UUID, number int, status string, accrual float64, uploadedAt time.Time) (*Order, error) {
	orderNumber, err := NewOrderNumber(number)
	if err != nil {
		return nil, NewOrderError(id.String(), err)
	}

	orderStatus, err := NewOrderStatus(status)
	if err != nil {
		return nil, NewOrderError(id.String(), err)
	}

	return &Order{
		id:         id,
		userID:     userID,
		number:     *orderNumber,
		accrual:    accrual,
		uploadedAt: uploadedAt,
		status:     *orderStatus,
	}, nil
}
