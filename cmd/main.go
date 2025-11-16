package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/milyrock/PR-Reviewer/internal/app"
	"github.com/milyrock/PR-Reviewer/internal/config"
	v1 "github.com/milyrock/PR-Reviewer/internal/handlers/v1"
	"github.com/milyrock/PR-Reviewer/internal/repository"
)

func main() {
	cfg, err := config.ReadConfig("./config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	db, err := app.InitDB(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to init db: %v", err)
	}
	defer db.Close()

	repo := repository.NewRepository(db)

	r := mux.NewRouter()

	api := v1.NewAPI(repo)
	api.RegisterHandlers(r)

	log.Println("Server starting on port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
