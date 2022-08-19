package app

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/sreway/gophermart/internal/controller/listener"
	"github.com/sreway/gophermart/internal/controller/processing"

	"github.com/sreway/gophermart/internal/usecase/accrual"
	"github.com/sreway/gophermart/pkg/httpclient"
	"github.com/sreway/gophermart/pkg/kafkaclient"

	"github.com/go-chi/chi/v5"

	"github.com/sreway/gophermart/config"
	"github.com/sreway/gophermart/internal/controller/http"
	"github.com/sreway/gophermart/internal/usecase/balance"
	"github.com/sreway/gophermart/internal/usecase/order"
	"github.com/sreway/gophermart/internal/usecase/queue"
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
	newOrdersChan := make(chan string)

	hc := httpclient.New(httpclient.WithBaseURL(cfg.Accrual.Address))

	pg, err := postgres.New(ctx, cfg.Postgres.DSN)
	if err != nil {
		logger.Fatal(err)
	}

	pgl, err := listener.NewPGListener(ctx, pg, cfg.Postgres.ListenChannel)
	if err != nil {
		logger.Fatal(err)
	}

	producer, err := kafkaclient.NewProducer(ctx, cfg.Kafka.BrokerNetwork, cfg.Kafka.BrokerAddress,
		cfg.Kafka.Topic, cfg.Kafka.Partition)
	if err != nil {
		logger.Fatal(err)
	}

	consumer := kafkaclient.NewConsumer(cfg.Kafka.BrokerAddress, cfg.Kafka.Topic, cfg.Kafka.GroupID)

	or := repo.NewOrderRepo(pg)
	ar := repo.NewAccrualRepo(hc)
	ur := repo.NewUserRepo(pg)
	br := repo.NewBalanceRepo(pg)
	wr := repo.NewWithdraw(pg)
	qr := repo.NewQueueRepo(producer, consumer)

	uc := user.New(ur)
	oc := order.New(or)
	bc := balance.New(br)
	wc := withdraw.New(wr)
	ac := accrual.New(ar, or)
	qc := queue.New(qr)

	pc, err := processing.NewOrders(qc, ac)
	if err != nil {
		logger.Fatal(err)
	}

	wg := new(sync.WaitGroup)
	wg.Add(3)

	// run http server
	go func(cfg *config.Config, stop chan os.Signal) {
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
	}(cfg, systemSignals)

	// run listen new order from postgres
	go func() {
		defer wg.Done()
		pgl.Listen(ctx, newOrdersChan)
	}()

	// send new order in kafka queue
	go func() {
		defer wg.Done()
		pc.SendToQueue(ctx, newOrdersChan)
	}()

	// processing new order in accrual service
	go func() {
		defer wg.Done()
		pc.CheckInAccrual(ctx)
	}()

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
	pgl.Release()
	pg.Close()
	producer.Close()
	consumer.Close()
	os.Exit(exitCode)
}
