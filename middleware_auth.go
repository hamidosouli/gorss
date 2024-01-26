package main

import (
	"fmt"
	"github.com/hamidosouli/rssaggregator/internal/auth"
	"github.com/hamidosouli/rssaggregator/internal/database"
	"net/http"
)

type authHandler func(w http.ResponseWriter, r *http.Request, user database.User)

func (cfg *apiConfig) middlewareAuth(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetApiKey(r.Header)
		if err != nil {
			responseError(w, 401, fmt.Sprintf("Auth error: %v", err))
			return
		}
		user, err := cfg.DB.GetUserByApiKey(r.Context(), apiKey)
		if err != nil {
			responseError(w, 403, fmt.Sprintf("Invalid apiKey: %v", err))
			return
		}
		handler(w, r, user)
	}
}
