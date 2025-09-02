package auth

import (
	"net/http"
	"strings"
	"time"

	"github.com/Arpitmovers/reviewservice/internal/config"
	logger "github.com/Arpitmovers/reviewservice/internal/logging"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateJWT creates a signed token
func GenerateJWT(username string, cfg *config.Config) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour)

	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	str, err := token.SignedString([]byte(cfg.JwtKey))
	if err != nil {
		logger.Logger.Error("failed to sign JWT", zap.String("username", username), zap.Error(err))
		return "", err
	}

	return str, nil
}

func JWTAuthMiddleware(cfg *config.Config, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			logger.Logger.Warn("missing Authorization header", zap.String("path", r.URL.Path), zap.String("method", r.Method))
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			logger.Logger.Warn("invalid Authorization header format", zap.String("header", authHeader), zap.String("path", r.URL.Path))
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		tokenStr := parts[1]
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JwtKey), nil
		})

		if err != nil || !token.Valid {
			logger.Logger.Warn("invalid or expired JWT token", zap.String("path", r.URL.Path), zap.Error(err))
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
