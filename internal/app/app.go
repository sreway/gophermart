package app

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

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

	newOrdersChannel := make(chan string)
	wg := new(sync.WaitGroup)
	wg.Add(3)

	// run http server
	go func(cfg *config.Config, pg *postgres.Postgres, stop chan os.Signal) {
		ur := repo.NewUserRepo(pg)
		or := repo.NewOrderRepo(pg)
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
	go orderListner(ctx, wg, pg, cfg, newOrdersChannel)

	// store order in kafka queue
	go storeOrderInKafka(ctx, wg, cfg, newOrdersChannel)

	// processing orders
	go processingOrder(ctx, wg, pg, cfg)

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
	os.Exit(exitCode)
}
