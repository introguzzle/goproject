package auth

import (
	"github.com/golang-jwt/jwt/v4"
	"goproject/src/env"
)

var Secret = []byte(env.Get("JWT_SECRET_KEY").Value)

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}
