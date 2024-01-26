package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/hamidosouli/rssaggregator/internal/database"
	"net/http"
	"time"
)

func (apiCfg *apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		responseError(w, 400, fmt.Sprintf("Error parsing json: %v", err))
		return
	}
	feed, err := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.URL,
		UserID:    user.ID,
	})
	if err != nil {
		responseError(w, 500, fmt.Sprintf("Error creating user: %v", err))
		return
	}

	responseJson(w, 201, fromDBToFeed(feed))
}

func (apiCfg *apiConfig) handlerGetFeed(w http.ResponseWriter, r *http.Request) {
	feeds, err := apiCfg.DB.GetFeeds(r.Context())
	if err != nil {
		responseError(w, 500, fmt.Sprintf("Error creating user: %v", err))
		return
	}
	responseJson(w, 201, fromDBToFeeds(feeds))
}
