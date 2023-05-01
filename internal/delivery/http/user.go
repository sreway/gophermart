package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/exp/slog"

	"github.com/sreway/gophermart/internal/domain"
)

func (d *delivery) userRegister(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	creds := new(credentialsRequest)

	if err := decoder.Decode(creds); err != nil {
		d.logger.Error("failed decode request body",
			slog.String("handler", "userRegister"),
			slog.Any("err", err))
		handelUserErr(w, ErrDecodeBody)
		return
	}

	user, err := d.user.Register(r.Context(), creds.Login, creds.Password)
	if err != nil {
		d.logger.Error("failed register user",
			slog.Any("err", err),
			slog.String("user_login", creds.Login),
		)
		handelUserErr(w, err)
		return
	}

	err = d.setAuthCookie(w, user.ID().String(), user.Login())
	if err != nil {
		handelUserErr(w, err)
	}

	w.WriteHeader(http.StatusOK)
}

func (d *delivery) userLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	creds := new(credentialsRequest)

	if err := decoder.Decode(creds); err != nil {
		d.logger.Error("failed decode request body",
			slog.String("handler", "userLogin"),
			slog.Any("err", err))
		handelUserErr(w, ErrDecodeBody)
		return
	}

	user, err := d.user.Login(r.Context(), creds.Login, creds.Password)
	if err != nil {
		handelUserErr(w, err)
		return
	}

	err = d.setAuthCookie(w, user.ID().String(), user.Login())
	if err != nil {
		handelUserErr(w, err)
	}

	w.WriteHeader(http.StatusOK)
}

func (d *delivery) userLogOut(w http.ResponseWriter, r *http.Request) {
	cookie := new(http.Cookie)
	cookie.Name = d.config.Auth.CookieName
	cookie.Value = ""
	cookie.MaxAge = 0
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
}

func (d *delivery) setAuthCookie(w http.ResponseWriter, userID, userLogin string) error {
	expirationTime := &jwt.NumericDate{Time: time.Now().Add(d.config.Auth.TokenTTL)}
	claim := claims{
		UserID:    userID,
		UserLogin: userLogin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: expirationTime,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString([]byte(d.config.Auth.Key))
	if err != nil {
		d.logger.Error("failed signed jwt token", slog.Any("error", err))
		return ErrSignedJWT
	}

	cookie := new(http.Cookie)
	cookie.Name = d.config.Auth.CookieName
	cookie.Value = tokenString
	cookie.Expires = expirationTime.Time

	http.SetCookie(w, cookie)
	return nil
}

func handelUserErr(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrDecodeBody):
		w.WriteHeader(http.StatusBadRequest)
	case errors.Is(err, domain.ErrIncorrectData):
		w.WriteHeader(http.StatusBadRequest)
	case errors.Is(err, domain.ErrAlreadyExist):
		w.WriteHeader(http.StatusConflict)
	case errors.Is(err, ErrSignedJWT):
		w.WriteHeader(http.StatusInternalServerError)
	case errors.Is(err, domain.ErrNotFound):
		w.WriteHeader(http.StatusNotFound)
	case errors.Is(err, domain.ErrUserUnauthorized):
		w.WriteHeader(http.StatusUnauthorized)
	default:
		w.WriteHeader(http.StatusNotImplemented)
	}
}
