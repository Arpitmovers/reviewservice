package auth

import (
	"encoding/json"
	"net/http"

	"github.com/Arpitmovers/reviewservice/internal/config"
	logger "github.com/Arpitmovers/reviewservice/internal/logging"
	"go.uber.org/zap"
)

type Creds struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Logger.Error("failed to write JSON response", zap.Error(err))
	}
}

func LoginHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var creds Creds

		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			logger.Logger.Error("invalid login request", zap.Error(err))
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request"})
			return
		}

		if creds.Username != cfg.ApiUser || creds.Password != cfg.ApiPwd {
			logger.Logger.Warn("invalid credentials attempt", zap.String("username", creds.Username))
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
			return
		}

		token, err := GenerateJWT(creds.Username, cfg)
		if err != nil {
			logger.Logger.Error("failed to generate JWT", zap.String("username", creds.Username), zap.Error(err))
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})
			return
		}

		writeJSON(w, http.StatusOK, map[string]string{"token": token})
	}
}
