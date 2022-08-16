package app

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/sreway/gophermart/config"
	"github.com/sreway/gophermart/internal/controller/listener"
	"github.com/sreway/gophermart/internal/entity"
	"github.com/sreway/gophermart/internal/usecase/accrual"
	"github.com/sreway/gophermart/internal/usecase/repo"
	"github.com/sreway/gophermart/pkg/httpclient"
	"github.com/sreway/gophermart/pkg/kafkaclient"
	"github.com/sreway/gophermart/pkg/logger"
	"github.com/sreway/gophermart/pkg/postgres"
)

type Stack struct {
	Values []string
	lock   *sync.Mutex
}

func (s *Stack) Put(d string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.Values = append(s.Values, d)
}

func (s *Stack) Get() *string {
	if len(s.Values) > 0 {
		s.lock.Lock()
		defer s.lock.Unlock()
		d := s.Values[0]
		s.Values = s.Values[1:]
		return &d
	}
	return nil
}

func NewOrderStack() *Stack {
	return &Stack{make([]string, 0), &sync.Mutex{}}
}

func orderListner(ctx context.Context, wg *sync.WaitGroup, pg *postgres.Postgres, cfg *config.Config, data chan<- string) {
	defer func() {
		logger.Info("Stop orderListner")
		wg.Done()
	}()
	pgl, err := listener.NewPgListner(ctx, pg, cfg.Postgres.ListenChannel)
	if err != nil {
		logger.Fatal(err)
	}
	defer pgl.Release()

	pgl.Listen(ctx, data)
}

func storeOrderInKafka(ctx context.Context, wg *sync.WaitGroup, cfg *config.Config, data <-chan string) {
	defer func() {
		logger.Info("Stop storeOrderInKafka")
		wg.Done()
	}()

	producer, err := kafkaclient.NewProducer(ctx, cfg.Kafka.BrokerNetwork, cfg.Kafka.BrokerAddress,
		cfg.Kafka.Topic, cfg.Kafka.Partition)
	if err != nil {
		logger.Fatal(err)
	}
	defer producer.Close()

	for {
		select {
		case <-ctx.Done():
			return
		case orderNumber := <-data:
			err = producer.Write(orderNumber)
			if err != nil {
				logger.Error(err)
			}
		}
	}
}

func processingOrder(ctx context.Context, wg *sync.WaitGroup, pg *postgres.Postgres, cfg *config.Config) {
	defer func() {
		logger.Info("Stop processingOrder")
		wg.Done()
	}()

	hc := httpclient.New(httpclient.WithBaseURL(cfg.Accrual.Address))
	or := repo.NewOrderRepo(pg)
	ar := repo.NewAccrualRepo(hc.Client)
	ac := accrual.New(ar, or)

	consumer := kafkaclient.NewConsumer(cfg.Kafka.BrokerAddress, cfg.Kafka.Topic, cfg.Kafka.GroupID)
	defer consumer.Close()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := consumer.Read(ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					return
				}
				logger.Error(err)
				continue
			}

			if msg == nil {
				continue
			}

			a, err := ac.Get(ctx, string(msg.Value))
			if err != nil {
				var rateLimitedError *entity.ErrRateLimited
				var httpClientError *entity.ErrHTTPClient

				switch {
				case errors.As(err, &rateLimitedError):
					logger.Error(rateLimitedError)
					time.Sleep(rateLimitedError.RetryAfter)
					continue
				case errors.As(err, &httpClientError):
					logger.Error(httpClientError)
					continue
				}
			}

			err = ac.UpdateOrderStatus(ctx, a)
			if err != nil {
				logger.Error(err)
				continue
			}

			switch entity.OrderStatus(a.Status) {
			case entity.OrderStatusInvalid, entity.OrdertStatusProcessed:
				logger.Infof("Success processed order %s", string(msg.Value))
			default:
				err = consumer.Commit(ctx, msg)
				if err != nil {
					logger.Error(err)
					continue
				}
			}
		}
	}
}
