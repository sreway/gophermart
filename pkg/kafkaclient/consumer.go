package kafkaclient

import (
	"context"

	"github.com/segmentio/kafka-go"

	"github.com/sreway/gophermart/pkg/logger"
)

type Consumer struct {
	r *kafka.Reader
}

func NewConsumer(brokerAddress string, topic string, groupid string) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{brokerAddress},
		Topic:   topic,
		GroupID: groupid,
	})
	logger.Info("NewConsumer: success connect")
	return &Consumer{
		r: r,
	}
}

func (c *Consumer) Read(ctx context.Context) (*kafka.Message, error) {
	msg, err := c.r.ReadMessage(ctx)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

func (c *Consumer) Fetch(ctx context.Context) (*kafka.Message, error) {
	msg, err := c.r.FetchMessage(ctx)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

func (c *Consumer) Commit(ctx context.Context, msg *kafka.Message) error {
	err := c.r.CommitMessages(ctx, *msg)
	if err != nil {
		return err
	}

	return nil
}

func (c *Consumer) Close() {
	err := c.r.Close()
	if err != nil {
		logger.Fatal(err)
	}
}
