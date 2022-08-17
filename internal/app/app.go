package app

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/sreway/gophermart/internal/usecase/accrual"
	"github.com/sreway/gophermart/pkg/httpclient"
	"github.com/sreway/gophermart/pkg/kafkaclient"

	"github.com/go-chi/chi/v5"

	"github.com/sreway/gophermart/config"
	"github.com/sreway/gophermart/internal/controller/http"
	"github.com/sreway/gophermart/internal/usecase/balance"
	"github.com/sreway/gophermart/internal/usecase/order"
	"github.com/sreway/gophermart/internal/usecase/repo"
	"github.com/sreway/gophermart/internal/usecase/user"
	"github.com/sreway/gophermart/internal/usecase/withdraw"
	"github.com/sreway/gophermart/pkg/auth"
	"github.com/sreway/gophermart/pkg/httpserver"
	"github.com/sreway/gophermart/pkg/logger"
	"github.com/sreway/gophermart/pkg/postgres"
)

func Run(cfg *config.Config) {
	ctx, cancel := context.WithCancel(context.Background())
	systemSignals := make(chan os.Signal, 1)
	signal.Notify(systemSignals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	exitChan := make(chan int)

	pg, err := postgres.New(ctx, cfg.Postgres.DSN)
	if err != nil {
		logger.Fatal(err)
	}

	producer, err := kafkaclient.NewProducer(ctx, cfg.Kafka.BrokerNetwork, cfg.Kafka.BrokerAddress,
		cfg.Kafka.Topic, cfg.Kafka.Partition)
	if err != nil {
		logger.Fatal(err)
	}

	consumer := kafkaclient.NewConsumer(cfg.Kafka.BrokerAddress, cfg.Kafka.Topic, cfg.Kafka.GroupID)

	newOrdersChannel := make(chan string)
	hc := httpclient.New(httpclient.WithBaseURL(cfg.Accrual.Address))

	oq := repo.NewQueueRepo(producer, consumer)
	or := repo.NewOrderRepo(pg)
	ar := repo.NewAccrualRepo(hc.Client)
	ac := accrual.New(ar, or)

	wg := new(sync.WaitGroup)
	wg.Add(3)

	// run http server
	go func(cfg *config.Config, pg *postgres.Postgres, stop chan os.Signal) {
		ur := repo.NewUserRepo(pg)
		br := repo.NewBalanceRepo(pg)
		wr := repo.NewWithdraw(pg)
		uc := user.New(ur)
		oc := order.New(or)
		bc := balance.New(br)
		wc := withdraw.New(wr)
		r := chi.NewRouter()
		jwt := auth.New(cfg.Server.Auth.JWT.Key, cfg.Server.Auth.JWT.TokenTTL)
		http.NewRouter(r, cfg, jwt, uc, oc, bc, wc)
		httpServer := httpserver.New(r, httpserver.Addr(cfg.Server.HTTP.Address))
		err = httpServer.ListenAndServe()
		if err != nil {
			logger.Errorf("httpserver: %v\n", err)
			stop <- syscall.SIGSTOP
			return
		}
	}(cfg, pg, systemSignals)

	// run listen new order from postgres
	go orderListener(ctx, wg, pg, cfg, newOrdersChannel)

	// store order in kafka queue
	go storeOrderQueue(ctx, wg, oq, newOrdersChannel)

	// processing orders
	go processingOrder(ctx, wg, ac, oq)

	// listen system signals for graceful shutdown
	go func() {
		for {
			systemSignal := <-systemSignals
			switch systemSignal {
			case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				logger.Info("signal triggered")
				exitChan <- 0
			default:
				logger.Warn("unknown signal")
				exitChan <- 1
			}
		}
	}()

	exitCode := <-exitChan
	cancel()
	wg.Wait()
	pg.Close()
	producer.Close()
	consumer.Close()
	os.Exit(exitCode)
}
