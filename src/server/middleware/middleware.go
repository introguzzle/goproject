package middleware

import (
	"context"
	"encoding/json"
	"github.com/golang-jwt/jwt/v4"
	_ "github.com/golang-jwt/jwt/v4"
	"goproject/src/server/auth"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	STATUS string = "$STATUS"
)

type Handler func(w http.ResponseWriter, r *http.Request)

type Middleware interface {
	Handle(handler Handler) Handler
	Ignore() []string
}

type RestMiddleware struct{}

func (jm *RestMiddleware) Handle(handler Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("X-Frame-Options", "DENY")

		handler(w, r)
	}
}

func (jm *RestMiddleware) Ignore() []string {
	return []string{"*"}
}

type LoggingMiddleware struct{}

func (lm *LoggingMiddleware) Handle(handler Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf(
			"Request accepted: %s %s",
			r.Method,
			r.URL.Path,
		)

		handler(w, r)
		end := time.Since(start)

		log.Printf(
			"Response sent: %s %s, Status: %s, Duration: %s",
			r.Method,
			r.RequestURI,
			w.Header().Get(STATUS),
			end,
		)

		w.Header().Del(STATUS)
	}
}

func (lm *LoggingMiddleware) Ignore() []string {
	return []string{"*"}
}

type AuthenticationMiddleware struct{}

func invalidToken(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	_, _ = w.Write(wrap("Invalid token", http.StatusUnauthorized))
}

func wrap(s string, code int) []byte {
	j, err := json.Marshal(map[string]string{
		"error":  s,
		"status": strconv.Itoa(code),
	})

	if err != nil {
		return nil
	}

	return j
}

func (am *AuthenticationMiddleware) Handle(handler Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, ignore := range am.Ignore() {
			if strings.HasPrefix(r.URL.Path, ignore) {
				handler(w, r)
				return
			}
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write(wrap("Authorization header missing", http.StatusUnauthorized))
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.ParseWithClaims(tokenString, &auth.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return auth.Secret, nil
		})

		if err != nil {
			invalidToken(w)
			return
		}

		claims, ok := token.Claims.(*auth.Claims)
		if !ok || !token.Valid {
			invalidToken(w)
			return
		}

		ctx := context.WithValue(r.Context(), "username", claims.Username)
		r = r.WithContext(ctx)

		handler(w, r)
	}
}

func (am *AuthenticationMiddleware) Ignore() []string {
	return []string{"/api/v1/auth"}
}

type ValidationMiddleware struct{}

func (vm *ValidationMiddleware) Handle(Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
