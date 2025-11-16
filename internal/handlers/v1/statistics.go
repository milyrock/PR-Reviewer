package v1

import (
	"encoding/json"
	"net/http"

	"github.com/milyrock/PR-Reviewer/internal/repository"
	"github.com/milyrock/PR-Reviewer/internal/service"
)

type StatisticsHandler struct {
	service *service.StatisticsService
}

func NewStatisticsHandler(repo *repository.Repository) *StatisticsHandler {
	return &StatisticsHandler{service: service.NewStatisticsService(repo)}
}

func (h *StatisticsHandler) GetStatistics(w http.ResponseWriter, r *http.Request) {
	stats, err := h.service.GetStatistics()
	if err != nil {
		writeError(w, statusInternalError, errorCodeInternalError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
