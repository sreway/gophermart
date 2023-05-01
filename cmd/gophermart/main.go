package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"sync"

	"golang.org/x/exp/slog"

	"github.com/sreway/gophermart/internal/config"
	"github.com/sreway/gophermart/internal/delivery/http"
	accrual "github.com/sreway/gophermart/internal/repository/accrual/http"
	"github.com/sreway/gophermart/internal/repository/storage/postgres"
	"github.com/sreway/gophermart/internal/usecases/balance"
	"github.com/sreway/gophermart/internal/usecases/order"
	"github.com/sreway/gophermart/internal/usecases/user"
	"github.com/sreway/gophermart/internal/usecases/withdraw"
)

var (
	serverAddress  string
	dsn            string
	accrualAddress string
	lookupEnv      = []string{
		"RUN_ADDRESS", "DATABASE_URI", "ACCRUAL_SYSTEM_ADDRESS",
	}
)

func init() {
	flag.StringVar(
		&serverAddress, "a", serverAddress, "Server address: host:port")
	flag.StringVar(
		&dsn, "d", dsn, "PostgreSQL data source name")
	flag.StringVar(
		&accrualAddress, "r", accrualAddress, "Accrual system address")
	flag.Parse()

	for _, env := range lookupEnv {
		_, exist := os.LookupEnv(env)
		if exist {
			continue
		}

		switch env {
		case "RUN_ADDRESS":
			_ = os.Setenv(env, serverAddress)
		case "DATABASE_URI":
			_ = os.Setenv(env, dsn)
		case "ACCRUAL_SYSTEM_ADDRESS":
			_ = os.Setenv(env, accrualAddress)
		}
	}
}

func main() {
	var code int

	log := slog.New(slog.NewJSONHandler(os.Stdout).
		WithAttrs([]slog.Attr{slog.String("service", "gophermart")}),
	)

	exit := make(chan int)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer func() {
		stop()
		os.Exit(code)
	}()

	wg := new(sync.WaitGroup)
	wg.Add(2)

	go func() {
		defer wg.Done()

		cfg, err := config.New()
		if err != nil {
			log.Error("failed initialize config", slog.Any("err", err))
			exit <- 1
			return
		}
		storage, err := postgres.New(ctx, cfg.Storage())
		if err != nil {
			log.Error("failed initialize postgres repository", slog.Any("err", err))
			exit <- 1
			return
		}
		accrualRepo := accrual.New(cfg.Accrual())

		userService := user.New(storage)
		orderService := order.New(cfg.Orders(), storage, accrualRepo)
		balanceService := balance.New(storage)
		withdrawService := withdraw.New(storage)

		go func() {
			defer wg.Done()
			err = orderService.ProcNewOrder(ctx, cfg.Orders())
			if err != nil {
				log.Error("failed processing new order", slog.Any("err", err))
				exit <- 1
				return
			}
		}()

		httpServer := http.New(cfg.HTTP(), userService, orderService, balanceService, withdrawService)
		err = httpServer.Run(ctx)
		if err != nil {
			log.Error("failed run http server", slog.Any("err", err))
			exit <- 1
			return
		}
	}()

	go func() {
		<-ctx.Done()
		stop()
		exit <- 0
		log.Info("trigger graceful shutdown app")
	}()

	code = <-exit
	wg.Done()
}
