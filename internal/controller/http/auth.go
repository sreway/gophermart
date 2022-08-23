package http

import (
	"context"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/jwt"
)

type contextKey string

const (
	userIDKey    contextKey = "user_id"
	userLoginKey contextKey = "user_login"
)

func (c contextKey) String() string {
	return string(c)
}

func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, claims, err := jwtauth.FromContext(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if token == nil || jwt.Validate(token) != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		userID, ok := claims[userIDKey.String()]
		if !ok {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}

		userLogin, ok := claims[userLoginKey.String()]
		if !ok {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}

		ctx := context.WithValue(r.Context(), userLoginKey, userLogin)
		ctx = context.WithValue(ctx, userIDKey, userID)

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
