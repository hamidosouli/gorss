package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func responseJson(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("failed to matshal json response %v", payload)
		w.WriteHeader(500)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(data)
	if err != nil {
		return
	}
}

func responseError(w http.ResponseWriter, code int, err string) {
	if code >= 500 {
		log.Println("responding with 5XX level err:", err)
	}

	type errResponse struct {
		Error string `json:"error"`
	}
	responseJson(w, code, errResponse{
		Error: err,
	})
}
