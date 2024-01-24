package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
		return
	}
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in env")
	}

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/error", handleErr)

	router.Mount("/v1", v1Router)

	server := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("server starting on port %v", portString)

	listenErr := server.ListenAndServe()
	if listenErr != nil {
		log.Fatal(listenErr)
		return
	}

	fmt.Println("PORT is:", portString)
}
