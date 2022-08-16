package entity

import (
	"encoding/json"
	"time"

	"github.com/go-playground/validator/v10"
)

type OrderStatus string

const (
	OrderStatusNew        OrderStatus = "NEW"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusInvalid    OrderStatus = "INVALID"
	OrdertStatusProcessed OrderStatus = "PROCESSED"
)

func AllowedOrderStatus() []OrderStatus {
	return []OrderStatus{OrderStatusNew, OrderStatusProcessing, OrderStatusInvalid, OrdertStatusProcessed}
}

func (os OrderStatus) String() string {
	for _, v := range AllowedOrderStatus() {
		if v == os {
			return string(v)
		}
	}
	return "Invalid status name"
}

func (os OrderStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(os.String())
}

type Order struct {
	ID         uint        `json:"-"`
	Number     string      `json:"number" db:"number" validate:"required"`
	UserID     uint        `json:"-" db:"user_id" validate:"required"`
	Status     OrderStatus `json:"status" db:"status"`
	Accrual    float64     `json:"accrual,omitempty" db:"accrual"`
	UploadedAt time.Time   `json:"uploaded_at" db:"uploaded_at"`
}

func (order *Order) Validate() error {
	validate := validator.New()
	err := validate.Struct(order)
	if err != nil {
		return err
	}
	return nil
}

func (order Order) MarshalJSON() ([]byte, error) {
	type OrderAlias Order
	aliasValue := struct {
		OrderAlias
		UploadedAt string `json:"uploaded_at"`
	}{
		OrderAlias: OrderAlias(order),
		UploadedAt: order.UploadedAt.Format(time.RFC3339),
	}
	return json.Marshal(aliasValue)
}
