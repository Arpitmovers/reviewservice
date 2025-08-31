package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Arpitmovers/reviewservice/internal/config"
)

type Creds struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoginHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var creds Creds

		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		if creds.Username != cfg.ApiUser || creds.Password != cfg.ApiPwd {

			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		token, err := GenerateJWT(creds.Username, cfg)
		if err != nil {
			fmt.Println("failed to generate token", err)
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"token": token})
	}
}
