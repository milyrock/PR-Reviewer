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

	teamHandler := v1.NewTeamHandler(repo)
	userHandler := v1.NewUserHandler(repo)
	prHandler := v1.NewPRHandler(repo)
	statisticsHandler := v1.NewStatisticsHandler(repo)

	r := mux.NewRouter()

	r.HandleFunc("/health", v1.Health).Methods("GET")

	r.HandleFunc("/team/add", teamHandler.AddTeam).Methods("POST")
	r.HandleFunc("/team/get", teamHandler.GetTeam).Methods("GET")

	r.HandleFunc("/users/setIsActive", userHandler.SetIsActive).Methods("POST")
	r.HandleFunc("/users/getReview", userHandler.GetReview).Methods("GET")

	r.HandleFunc("/pullRequest/create", prHandler.CreatePR).Methods("POST")
	r.HandleFunc("/pullRequest/merge", prHandler.MergePR).Methods("POST")
	r.HandleFunc("/pullRequest/reassign", prHandler.ReassignPR).Methods("POST")

	r.HandleFunc("/statistics", statisticsHandler.GetStatistics).Methods("GET")

	log.Println("Server starting on port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
