package v1

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/milyrock/PR-Reviewer/internal/models"
	"github.com/milyrock/PR-Reviewer/internal/repository"
	"github.com/milyrock/PR-Reviewer/internal/service"
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
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"pr": createdPR,
	}); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
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
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"pr": pr,
	}); err != nil{
		log.Printf("failed to encode response: %v", err)
	}
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
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"pr":          updatedPR,
		"replaced_by": newReviewerID,
	}); err != nil{
		log.Printf("failed to encode response: %v", err)
	}
}
