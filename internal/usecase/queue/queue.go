package queue

import (
	"context"

	"github.com/sreway/gophermart/internal/entity"
	"github.com/sreway/gophermart/internal/usecase"
)

type Queue struct {
	repo usecase.QueueRepo
}

func New(queue usecase.QueueRepo) *Queue {
	return &Queue{
		repo: queue,
	}
}

func (q *Queue) Add(ctx context.Context, number string) error {
	return q.repo.Add(ctx, number)
}

func (q *Queue) Read(ctx context.Context) (*entity.QueueMsg, error) {
	return q.repo.Read(ctx)
}

func (q *Queue) Commit(ctx context.Context, msg *entity.QueueMsg) error {
	return q.repo.Commit(ctx, msg)
}
