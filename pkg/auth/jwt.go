package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type (
	JWT struct {
		SecretKey string
		TokenTTL  time.Duration
	}
)

func New(secretKey string, tokenTTL time.Duration) *JWT {
	return &JWT{SecretKey: secretKey, TokenTTL: tokenTTL}
}

func (j *JWT) NewToken(issuer string, claimsData map[string]interface{}) (string, error) {
	claims := jwt.MapClaims{
		"iss": issuer,
		"exp": time.Now().Add(j.TokenTTL).Unix(),
	}

	for k, v := range claimsData {
		claims[k] = v
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.SecretKey))
}

func (j *JWT) Exp(token string) (*time.Time, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("JWT_Exp unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.SecretKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("JWT_Exp: error get claims from token")
	}

	var expTime time.Time
	switch exp := claims["exp"].(type) {
	case float64:
		expTime = time.Unix(int64(exp), 0)
	default:
		return nil, fmt.Errorf("JWT_Exp: can't parse expires at")
	}

	return &expTime, nil
}
