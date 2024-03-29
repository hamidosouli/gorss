package main

import (
	"database/sql"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/hamidosouli/rssaggregator/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // include even though I am not calling it directly
	"log"
	"net/http"
	"os"
	"time"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {

	//rssFeed, err := urlToFeed("https://www.wagslane.dev/index.xml")
	//if err != nil {
	//	log.Fatal(err)
	//	return
	//}
	//fmt.Printf("got rssFeed from url and description for channel is %v", rssFeed.Channel.Description)
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
		return
	}
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in env")
	}
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not found in env")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("error in connecting to database: ", err)
	}

	apiCfg := apiConfig{
		DB: database.New(conn),
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
	v1Router.Post("/users", apiCfg.handlerCreateUser)
	v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
	v1Router.Get("/feeds", apiCfg.handlerGetFeed)
	v1Router.Post("/feed_follow", apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollow))
	v1Router.Get("/feed_follow", apiCfg.middlewareAuth(apiCfg.handlerGetFeedFollow))
	v1Router.Get("/posts", apiCfg.middlewareAuth(apiCfg.handlerGetPostsByUser))
	v1Router.Delete("/feed_follow/{id}", apiCfg.middlewareAuth(apiCfg.handlerDeleteFeedFollow))
	v1Router.Post("/me", apiCfg.middlewareAuth(apiCfg.handlerGetUserByApiKey))

	router.Mount("/v1", v1Router)

	server := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("server starting on port %v", portString)

	go startScraping(apiCfg.DB, 10, time.Minute)
	listenErr := server.ListenAndServe()
	if listenErr != nil {
		log.Fatal(listenErr)
		return
	}

	fmt.Println("PORT is:", portString)
}
