package v1

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/milyrock/PR-Reviewer/internal/models"
	"github.com/milyrock/PR-Reviewer/internal/repository"
	"github.com/milyrock/PR-Reviewer/internal/service"
)

type TeamHandler struct {
	service *service.TeamService
}

func NewTeamHandler(repo *repository.Repository) *TeamHandler {
	return &TeamHandler{service: service.NewTeamService(repo)}
}

func (h *TeamHandler) AddTeam(w http.ResponseWriter, r *http.Request) {
	var req models.CreateTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, statusBadRequest, errorCodeInvalidRequest, errorMsgInvalidRequestBody)
		return
	}

	team, err := h.service.AddTeam(req)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCreated)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"team": team,
	}); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

func (h *TeamHandler) GetTeam(w http.ResponseWriter, r *http.Request) {
	teamName := r.URL.Query().Get("team_name")
	if teamName == "" {
		writeError(w, statusBadRequest, errorCodeInvalidRequest, errorMsgTeamNameRequired)
		return
	}

	team, err := h.service.GetTeam(teamName)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(team); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}
