package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joeavaikath/rssagg/internal/database"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {

	fmt.Println("Hello World!")

	// Loads into environment
	godotenv.Load(".env")

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT not found in .env")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL not found in .env")
	}

	// Create database connection
	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}

	db := database.New(conn)
	// Create a new apiConfig struct
	apiCfg := apiConfig{
		DB: database.New(conn),
	}

	go startScraping(db, 10, time.Minute)

	fmt.Println("Port:", portString)

	log.Printf("\n Create the main router...")
	router := chi.NewRouter()

	log.Printf("\n Set the cors handler with cors options...")
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	log.Printf("\n Create a subrouter...")
	v1Router := chi.NewRouter()

	log.Printf("\n Create a http method with a pattern and a handler func...")
	v1Router.Get("/healthz", handlerReadiness)

	v1Router.Get("/err", handlerErr)

	v1Router.Post("/users", apiCfg.handlerCreateUser)
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUser))

	v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
	v1Router.Get("/feeds", apiCfg.handlerGetFeeds)

	v1Router.Post("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollow))
	v1Router.Get("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerGetFeedFollows))
	v1Router.Delete("/feed_follows/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.handlerDeleteFeedFollow))

	v1Router.Get("/posts", apiCfg.middlewareAuth(apiCfg.handlerGetPostsForUser))

	router.Mount("/v1", v1Router)

	log.Printf("\n Create a server(struct) with:\n\t a handler (router) \n\t an address...")
	srv := &http.Server{
		Handler: router,
		Addr:    "localhost:" + portString,
	}

	log.Printf("\n Server starting on port %v!", portString)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
