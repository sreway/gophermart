package kafkaclient

import (
	"context"

	"github.com/segmentio/kafka-go"

	"github.com/sreway/gophermart/pkg/logger"
)

type (
	Producer struct {
		conn *kafka.Conn
	}
)

func NewProducer(ctx context.Context, brokerNetwork, brokerAddress string, topic string, partition int) (*Producer, error) {
	conn, err := kafka.DialLeader(ctx, brokerNetwork, brokerAddress, topic, partition)
	if err != nil {
		logger.Fatalf("NewProducer: %v", err)
	}
	logger.Info("NewProducer: success connect")
	return &Producer{
		conn: conn,
	}, nil
}

func (p *Producer) Write(msg string) error {
	_, err := p.conn.WriteMessages(
		kafka.Message{Value: []byte(msg)})
	if err != nil {
		return err
	}
	return nil
}

func (p *Producer) Close() {
	err := p.conn.Close()
	if err != nil {
		logger.Fatalf("ProducerClose: %v", err)
	}
}
