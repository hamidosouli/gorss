package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/hamidosouli/rssaggregator/internal/database"
	"net/http"
	"time"
)

func (apiCfg *apiConfig) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedId uuid.UUID `json:"feed_id"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		responseError(w, 400, fmt.Sprintf("Error parsing json: %v", err))
		return
	}
	feedFollow, err := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedId,
	})
	if err != nil {
		responseError(w, 500, fmt.Sprintf("Error creating feed follows: %v", err))
		return
	}

	responseJson(w, 201, fromDBToFeedFollow(feedFollow))
}

func (apiCfg *apiConfig) handlerGetFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {

	feedFollows, err := apiCfg.DB.GetFeedFollows(r.Context(), user.ID)
	if err != nil {
		responseError(w, 500, fmt.Sprintf("Error getting feed follows: %v", err))
		return
	}

	responseJson(w, 201, fromDBToFeedsFollows(feedFollows))
}

func (apiCfg *apiConfig) handlerDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	id := chi.URLParam(r, "id")
	feedFollowId, err := uuid.Parse(id)
	if err != nil {
		responseError(w, 400, fmt.Sprintf("could not parse given uuid: %v", err))
		return
	}
	err = apiCfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     feedFollowId,
		UserID: user.ID,
	})
	if err != nil {
		responseError(w, 500, fmt.Sprintf("Error unfollowing feed: %v", err))
		return
	}

	responseJson(w, 204, struct{}{})
}
