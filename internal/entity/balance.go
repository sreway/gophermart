package entity

import (
	"encoding/json"
	"fmt"
)

type (
	Balance struct {
		ID        uint    `json:"-"`
		UserID    uint    `json:"-" db:"user_id" validate:"required"`
		Value     float64 `json:"current" db:"balance"`
		Withdrawn float64 `json:"withdrawn" db:"withdrawn"`
	}
	CustomFloat struct {
		Value  float64
		Digits int
	}
)

func (l CustomFloat) MarshalJSON() ([]byte, error) {
	s := fmt.Sprintf("%.*f", l.Digits, l.Value)
	return []byte(s), nil
}

func (b Balance) MarshalJSON() ([]byte, error) {
	type BalanceAlias Balance

	aliasValue := struct {
		BalanceAlias
		Value     CustomFloat `json:"current"`
		Withdrawn CustomFloat `json:"withdrawn"`
	}{
		BalanceAlias: BalanceAlias(b),
		Value:        CustomFloat{Value: b.Value, Digits: 2},
		Withdrawn:    CustomFloat{Value: b.Withdrawn, Digits: 2},
	}
	return json.Marshal(aliasValue)
}
