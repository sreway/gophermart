package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"

	"github.com/sreway/gophermart/config"
	"github.com/sreway/gophermart/internal/usecase"
	"github.com/sreway/gophermart/pkg/auth"
)

func NewRouter(r *chi.Mux, cfg *config.Config, auth *auth.JWT,
	u usecase.User, o usecase.Order, b usecase.Balance, w usecase.Withdraw,
) {
	jwtAuth := jwtauth.New("HS256", []byte(cfg.Server.Auth.JWT.Key), nil)
	ur := newUserRoutes(u, auth)
	or := newOrderRoutes(o)
	br := newBalanceRoutes(b)
	wr := newWithdrawRoutes(w)

	r.Use(middleware.Compress(cfg.Server.HTTP.CompressLevel, cfg.Server.HTTP.CompressTypes...))

	// public routes
	r.Group(func(r chi.Router) {
		r.Post("/api/user/register", ur.userRegister)
		r.Post("/api/user/login", ur.userLogin)
	})

	// private routes
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(jwtAuth))
		r.Use(Authenticator)

		// user routes
		r.Group(func(r chi.Router) {
			r.Get("/api/user/logout", ur.userLogOut)
		})

		// order routes
		r.Group(func(r chi.Router) {
			r.Post("/api/user/orders", or.orderAdd)
			r.Get("/api/user/orders", or.orderGet)
		})

		// balance routes
		r.Group(func(r chi.Router) {
			r.Get("/api/user/balance", br.balanceGet)
		})

		// withdraw routes
		r.Group(func(r chi.Router) {
			r.Post("/api/user/balance/withdraw", wr.withdrawAdd)
			r.Get("/api/user/withdrawals", wr.withdrawGet)
		})
	})
}
