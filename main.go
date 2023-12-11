package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/daunsid/rssagg/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	fmt.Println("Hello world")

	godotenv.Load(".env")

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("Port not found in environment")
	}

	dbString := os.Getenv("DB_URL")
	if dbString == "" {
		log.Fatal("db url not found in environment")
	}

	conn, err := sql.Open("postgres", dbString)
	if err != nil {
		log.Fatal("Can't connect to database")
	}

	queries := database.New(conn)
	if err != nil {
		log.Fatal("Can't connect to db: ", err)
	}

	apiCfg := apiConfig{
		DB: queries,
	}

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()

	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/err", handlerErr)
	v1Router.Post("/users", apiCfg.handlerCreateUser)
	router.Mount("/v1", v1Router)
	//router.HandleFunc("/ready", HandlerReadiness)

	srv := http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}
	log.Printf("Server stating on port %v", portString)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

	println("port: ", portString)

}
