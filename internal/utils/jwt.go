package utils

import (
	"context"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func SaveJWT(ctx context.Context, issuer string) error {
	jwt, err := GenerateJWT(issuer)

	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:     "jwt",
		Value:    jwt,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: 2,
	}

	err = grpc.SetHeader(ctx, metadata.Pairs("set-cookie", cookie.String()))

	if err != nil {
		return err
	}

	return nil
}

func GenerateJWT(issuer string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    issuer,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})

	jwt, err := claims.SignedString([]byte("secret-key"))

	if err != nil {
		return "", err
	}

	return jwt, nil
}

func VerifyJWT(j string) (*jwt.Token, *jwt.StandardClaims, error) {
	token, err := jwt.ParseWithClaims(j, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret-key"), nil
	})

	if err != nil {
		return nil, nil, err
	}

	claims := token.Claims.(*jwt.StandardClaims)

	return token, claims, err
}
