package domain

import (
	"github.com/google/uuid"
)

type Balance struct {
	id        uuid.UUID
	userID    uuid.UUID
	value     float64
	withdrawn float64
}

func (b *Balance) ID() uuid.UUID {
	return b.id
}

func (b *Balance) UserID() uuid.UUID {
	return b.userID
}

func (b *Balance) Value() float64 {
	return b.value
}

func (b *Balance) Withdrawn() float64 {
	return b.withdrawn
}

func NewBalance(userID uuid.UUID, value, withdrawn float64) *Balance {
	id := uuid.New()
	return &Balance{
		id:        id,
		userID:    userID,
		value:     value,
		withdrawn: withdrawn,
	}
}
