package http

import (
	"context"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
)

type (
	claims struct {
		UserID    string `json:"user_id"`
		UserLogin string `json:"user_login"`
		jwt.RegisteredClaims
	}
	ctxUserID    struct{}
	ctxUserLogin struct{}
)

func ValidateCookieToken(jwtSecretKey, jwtCookieName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(jwtCookieName)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			token := cookie.Value
			claim := claims{}

			parsedToken, err := jwt.ParseWithClaims(token, &claim, func(token *jwt.Token) (interface{}, error) {
				return []byte(jwtSecretKey), nil
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			if parsedToken == nil || !parsedToken.Valid {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), ctxUserID{}, claim.UserID)
			ctx = context.WithValue(ctx, ctxUserLogin{}, claim.UserLogin)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
