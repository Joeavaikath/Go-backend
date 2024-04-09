package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Hello World!")

	// Loads into environment
	godotenv.Load(".env")

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT not found in .env")
	}

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

	router.Mount("/v1", v1Router)

	log.Printf("\n Create a server(struct) with:\n\t a handler (router) \n\t an address...")
	srv := &http.Server{
		Handler: router,
		Addr:    "localhost:" + portString,
	}

	log.Printf("\n Server starting on port %v!", portString)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
