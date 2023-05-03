package http

import (
	"context"
	"errors"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"golang.org/x/exp/slog"

	"github.com/sreway/gophermart/internal/config"
	"github.com/sreway/gophermart/internal/usecases"
)

type (
	delivery struct {
		router   *chi.Mux
		config   *config.HTTP
		user     usecases.User
		order    usecases.Order
		balance  usecases.Balance
		withdraw usecases.Withdraw
		logger   *slog.Logger
	}
)

func New(config *config.HTTP, user usecases.User, order usecases.Order, balance usecases.Balance,
	withdraw usecases.Withdraw,
) *delivery {
	log := slog.New(slog.NewJSONHandler(os.Stdout).
		WithAttrs([]slog.Attr{slog.String("service", "http")}))
	d := &delivery{
		logger:   log,
		user:     user,
		order:    order,
		balance:  balance,
		withdraw: withdraw,
		config:   config,
	}
	return d
}

func (d *delivery) Run(ctx context.Context) error {
	d.router = d.initRouter()
	httpServer := &http.Server{
		Addr:    d.config.Address,
		Handler: d.router,
	}

	ctxServer, stopServer := context.WithCancel(context.Background())
	go func() {
		<-ctx.Done()
		d.logger.Info("trigger graceful shutdown http server")
		err := httpServer.Shutdown(ctxServer)
		if err != nil {
			d.logger.Error("shutdown http server", err)
		}
		stopServer()
	}()
	d.logger.Info("http service is ready to listen and serv")
	err := httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	<-ctxServer.Done()
	return nil
}
