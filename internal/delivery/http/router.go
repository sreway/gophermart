package http

import (
	"github.com/go-chi/chi/v5"
)

func (d *delivery) initRouter() *chi.Mux {
	router := chi.NewRouter()

	router.Route("/api", func(r chi.Router) {
		r.Route("/user", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Post("/register", d.userRegister)
				r.Post("/login", d.userLogin)
			})

			r.Group(func(r chi.Router) {
				r.Use(ValidateCookieToken(d.config.Auth.Key, d.config.Auth.CookieName))
				r.Get("/logout", d.userLogOut)
				r.Route("/orders", func(r chi.Router) {
					r.Post("/", d.orderAdd)
					r.Get("/", d.orderGet)
				})
				r.Route("/balance", func(r chi.Router) {
					r.Get("/", d.balanceGet)
					r.Post("/withdraw", d.withdrawAdd)
				})

				r.Route("/withdrawals", func(r chi.Router) {
					r.Get("/", d.withdrawGet)
				})
			})
		})
	})

	return router
}
