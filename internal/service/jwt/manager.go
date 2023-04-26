package jwt

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

type TokenManager interface {
	NewJWT(userId string, ttl time.Duration) (string, error)
	ParseAuthorizationHeader(accessToken string) (string, error)
	JWTMiddleware(next http.Handler) http.Handler
}

type claims struct {
	jwt.StandardClaims
	UserId string `json:"user_id"`
}

type manager struct {
	signingKey string
}

func NewManager(signingKey string) (*manager, error) {
	if signingKey == "" {
		return nil, errors.New("empty signing key")
	}
	return &manager{signingKey: signingKey}, nil
}

func (m *manager) NewJWT(userID string, ttl time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ttl).Unix(),
		},
		UserId: userID,
	})
	return token.SignedString([]byte(m.signingKey))
}

func (m *manager) ParseAuthorizationHeader(accessToken string) (string, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (i any, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.signingKey), nil
	})
	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("error get user claims from token")
	}

	return claims["sub"].(string), nil
}

func (m *manager) JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, err := m.checkToken(w, r)
		if err != nil {
			return
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *manager) checkToken(w http.ResponseWriter, r *http.Request) (context.Context, error) {
	authHeader := r.Header.Get("Authorization")
	bearerStr := strings.Split(authHeader, "Bearer ")

	if len(bearerStr) != 2 {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"parse auth header err"}`))
		return nil, errors.New("auth header bad")
	}
	parsedToken, err := jwt.Parse(bearerStr[1], func(jwtT *jwt.Token) (any, error) {
		if _, ok := jwtT.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("err signing %v", jwtT.Header["alg"])
		}
		return []byte(m.signingKey), nil
	})

	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"token parse err"}`))
		return nil, err
	}

	if !parsedToken.Valid {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"token not valid"}`))
		return nil, errors.New("token not valid")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"claims err"}`))
		return nil, errors.New("claims err")
	}
	ctx := context.WithValue(r.Context(), "userID", claims["user_id"])

	return ctx, nil
}
