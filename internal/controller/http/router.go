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

	r.Route("/api", func(r chi.Router) {
		r.Route("/user", func(r chi.Router) {
			// public
			r.Group(func(r chi.Router) {
				// user methods
				r.Post("/register", ur.userRegister)
				r.Post("/login", ur.userLogin)
			})
			// private
			r.Group(func(r chi.Router) {
				r.Use(jwtauth.Verifier(jwtAuth))
				r.Use(Authenticator)
				r.Get("/logout", ur.userLogOut)
				r.Post("/orders", or.orderAdd)
				r.Get("/orders", or.orderGet)

				r.Route("/balance", func(r chi.Router) {
					r.Get("/", br.balanceGet)
					r.Post("/withdraw", wr.withdrawAdd)
				})

				r.Route("/withdrawals", func(r chi.Router) {
					r.Get("/", wr.withdrawGet)
					r.Get("/v2", wr.withdrawGetPagination)
				})
			})
		})
	})
}
