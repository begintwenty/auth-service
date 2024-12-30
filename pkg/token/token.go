package token

import (
	"authservice/domain"
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

type TokenClaims struct {
	ID string `bson: "_id json:"id"`
	jwt.StandardClaims
}

func GenJWT(ctx context.Context, userId domain.UserID, rememberMe bool) (string, error) {
	var expiresAt int64
	switch rememberMe {
	case true:
		expiresAt = time.Now().Add((time.Minute * 60) * 24 * 30).Unix()
	case false:
		expiresAt = time.Now().Add((time.Minute * 60) * 24).Unix()
	}

	claims := TokenClaims{
		ID: userId.AsString(),
		StandardClaims: jwt.StandardClaims{
			IssuedAt: time.Now().Unix(),
			// One day if remember me is false, and 30 days if remember me is true.
			ExpiresAt: expiresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(viper.GetString("SECRET_KEY")))
}

func GenPasswordResetJWT(userId domain.UserID) (string, error) {
	expiresAt := time.Now().Add((time.Minute * 15)).Unix()
	claims := TokenClaims{
		ID: userId.AsString(),
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: expiresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(viper.GetString("SECRET_KEY")))
}

func VerifyJWT(tokenString string) (string, error) {
	claims := &TokenClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method in token")
		}
		return []byte(viper.GetString("SECRET_KEY")), nil

	})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", errors.New("Invalid  token")
	}

	return claims.ID, nil
}
