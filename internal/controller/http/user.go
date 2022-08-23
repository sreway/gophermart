package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/sreway/gophermart/pkg/logger"

	"github.com/sreway/gophermart/internal/entity"

	"github.com/sreway/gophermart/internal/usecase"
	"github.com/sreway/gophermart/pkg/auth"
)

type userRoutes struct {
	user usecase.User
	auth *auth.JWT
}

func newUserRoutes(u usecase.User, auth *auth.JWT) *userRoutes {
	return &userRoutes{u, auth}
}

func (ur *userRoutes) userRegister(w http.ResponseWriter, r *http.Request) {
	var u entity.User
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	if err := decoder.Decode(&u); err != nil {
		logger.Error(err)
		HandelErrUser(w, entity.NewUserError(u.Login, entity.ErrUserUnauthorized))
		return
	}

	err := ur.user.Register(r.Context(), &u)
	if err != nil {
		HandelErrUser(w, entity.NewUserError(u.Login, err))
		return
	}

	ts, err := ur.auth.NewToken(u.Login, map[string]interface{}{
		userLoginKey.String(): u.Login,
		userIDKey.String():    u.ID,
	})
	if err != nil {
		logger.Error(err)
		HandelErrUser(w, entity.NewUserError(u.Login, entity.ErrUserUnauthorized))
		return
	}

	exp, err := ur.auth.Exp(ts)
	if err != nil {
		HandelErrUser(w, entity.NewUserError(u.Login, entity.ErrUserUnauthorized))
		return
	}

	cookie := &http.Cookie{
		Name:     "jwt",
		Value:    ts,
		Path:     "/",
		Expires:  *exp,
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusOK)
}

func (ur *userRoutes) userLogin(w http.ResponseWriter, r *http.Request) {
	var u entity.User
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	if err := decoder.Decode(&u); err != nil {
		logger.Error(err)
		HandelErrUser(w, entity.NewUserError(u.Login, entity.ErrUserUnauthorized))
		return
	}

	user, err := ur.user.Login(r.Context(), &u)
	if err != nil {
		logger.Error(entity.NewUserError(u.Login, err))
		HandelErrUser(w, entity.NewUserError(u.Login, entity.ErrUserUnauthorized))
		return
	}

	ts, err := ur.auth.NewToken(user.Login, map[string]interface{}{
		userLoginKey.String(): user.Login,
		userIDKey.String():    user.ID,
	})
	if err != nil {
		logger.Error(err)
		HandelErrUser(w, entity.NewUserError(u.Login, entity.ErrUserUnauthorized))
		return
	}

	exp, err := ur.auth.Exp(ts)
	if err != nil {
		HandelErrUser(w, entity.NewUserError(u.Login, entity.ErrUserUnauthorized))
		return
	}

	cookie := &http.Cookie{
		Name:     "jwt",
		Value:    ts,
		Path:     "/",
		Expires:  *exp,
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
}

func (ur *userRoutes) userLogOut(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cookie := &http.Cookie{
		Name:     "jwt",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
}

func HandelErrUser(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, entity.ErrUserAlreadyExist):
		w.WriteHeader(http.StatusConflict)
	case errors.Is(err, entity.ErrUserIncorrectData):
		w.WriteHeader(http.StatusBadRequest)
	case errors.Is(err, entity.ErrUserUnauthorized):
		w.WriteHeader(http.StatusUnauthorized)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	logger.Error(err)
}
