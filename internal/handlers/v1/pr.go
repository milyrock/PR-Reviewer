package v1

import (
	"encoding/json"
	"net/http"

	"github.com/milyrock/PR-Reviewer/internal/service"
	"github.com/milyrock/PR-Reviewer/internal/models"
	"github.com/milyrock/PR-Reviewer/internal/repository"
)

type PRHandler struct {
	service *service.PRService
}

func NewPRHandler(repo *repository.Repository) *PRHandler {
	return &PRHandler{service: service.NewPRService(repo)}
}

func (h *PRHandler) CreatePR(w http.ResponseWriter, r *http.Request) {
	var req models.CreatePRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, statusBadRequest, errorCodeInvalidRequest, errorMsgInvalidRequestBody)
		return
	}

	createdPR, err := h.service.CreatePR(req)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"pr": createdPR,
	})
}

func (h *PRHandler) MergePR(w http.ResponseWriter, r *http.Request) {
	var req models.MergePRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, statusBadRequest, errorCodeInvalidRequest, errorMsgInvalidRequestBody)
		return
	}

	pr, err := h.service.MergePR(req)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"pr": pr,
	})
}

func (h *PRHandler) ReassignPR(w http.ResponseWriter, r *http.Request) {
	var req models.ReassignPRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, statusBadRequest, errorCodeInvalidRequest, errorMsgInvalidRequestBody)
		return
	}

	updatedPR, newReviewerID, err := h.service.ReassignPR(req)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"pr":          updatedPR,
		"replaced_by": newReviewerID,
	})
}
