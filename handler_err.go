package main

import "net/http"

func handleErr(w http.ResponseWriter, r *http.Request) {
	responseError(w, 400, "Something went wrong")
}
