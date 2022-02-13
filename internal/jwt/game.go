package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"time"
)

func NewGameJWT(channel string, uuid string, secret []byte) (string, error) {
	claims := &jwt.RegisteredClaims{
		Subject:  channel,
		IssuedAt: jwt.NewNumericDate(time.Now()),
		ID:       uuid,
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
}

func VerifyGameJWT(tokenString string, secret []byte) (*jwt.RegisteredClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return secret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		err := claims.Valid()
		if err != nil {
			return nil, err
		}

		return claims, nil
	}

	return nil, errors.New("token wasn't valid")
}
