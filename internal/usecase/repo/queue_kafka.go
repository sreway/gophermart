package repo

import (
	"context"

	"github.com/segmentio/kafka-go"

	"github.com/sreway/gophermart/internal/entity"
	"github.com/sreway/gophermart/pkg/kafkaclient"
)

type QueueRepo struct {
	producer *kafkaclient.Producer
	consumer *kafkaclient.Consumer
}

func NewQueueRepo(p *kafkaclient.Producer, c *kafkaclient.Consumer) *QueueRepo {
	return &QueueRepo{producer: p, consumer: c}
}

func (q *QueueRepo) Add(ctx context.Context, number string) error {
	_ = ctx
	return q.producer.Write(number)
}

func (q *QueueRepo) Read(ctx context.Context) (*entity.QueueMsg, error) {
	msg, err := q.consumer.Fetch(ctx)
	if err != nil {
		return nil, err
	}
	return &entity.QueueMsg{Msg: msg, Value: msg.Value}, err
}

func (q *QueueRepo) Commit(ctx context.Context, msg *entity.QueueMsg) error {
	return q.consumer.Commit(ctx, msg.Msg.(*kafka.Message))
}
