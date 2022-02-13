package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"time"
)

type FrontendJWT struct {
	Name string
	jwt.RegisteredClaims
}

func NewFrontendJWT(name string, id string, token string, issuer string, secret []byte) (string, error) {
	claims := &FrontendJWT{
		Name: name,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   id,
			Issuer:    issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 8)),
			ID:        token,
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
}

func VerifyFrontendJWT(tokenString string, secret []byte) (*FrontendJWT, error) {
	token, err := jwt.ParseWithClaims(tokenString, &FrontendJWT{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return secret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*FrontendJWT); ok && token.Valid {
		err := claims.Valid()
		if err != nil {
			return nil, err
		}

		return claims, nil
	}

	return nil, errors.New("token wasn't valid")
}
