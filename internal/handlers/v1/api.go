package v1

import (
	"github.com/gorilla/mux"
	"github.com/milyrock/PR-Reviewer/internal/repository"
)

type API struct {
	teamHandler       *TeamHandler
	userHandler       *UserHandler
	prHandler         *PRHandler
	statisticsHandler *StatisticsHandler
}

func NewAPI(repo *repository.Repository) *API {
	return &API{
		teamHandler:       NewTeamHandler(repo),
		userHandler:       NewUserHandler(repo),
		prHandler:         NewPRHandler(repo),
		statisticsHandler: NewStatisticsHandler(repo),
	}
}

func (a *API) RegisterHandlers(r *mux.Router) {
	a.registerHealthHandlers(r)
	a.registerTeamHandlers(r)
	a.registerUserHandlers(r)
	a.registerPRHandlers(r)
	a.registerStatisticsHandlers(r)
}

func (a *API) registerHealthHandlers(r *mux.Router) {
	r.HandleFunc("/health", Health).Methods("GET")
}

func (a *API) registerTeamHandlers(r *mux.Router) {
	r.HandleFunc("/team/add", a.teamHandler.AddTeam).Methods("POST")
	r.HandleFunc("/team/get", a.teamHandler.GetTeam).Methods("GET")
}

func (a *API) registerUserHandlers(r *mux.Router) {
	r.HandleFunc("/users/setIsActive", a.userHandler.SetIsActive).Methods("POST")
	r.HandleFunc("/users/getReview", a.userHandler.GetReview).Methods("GET")
}

func (a *API) registerPRHandlers(r *mux.Router) {
	r.HandleFunc("/pullRequest/create", a.prHandler.CreatePR).Methods("POST")
	r.HandleFunc("/pullRequest/merge", a.prHandler.MergePR).Methods("POST")
	r.HandleFunc("/pullRequest/reassign", a.prHandler.ReassignPR).Methods("POST")
}

func (a *API) registerStatisticsHandlers(r *mux.Router) {
	r.HandleFunc("/statistics", a.statisticsHandler.GetStatistics).Methods("GET")
}
