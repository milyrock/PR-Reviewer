package v1

import (
	"encoding/json"
	"net/http"

	"github.com/milyrock/PR-Reviewer/internal/service"
	"github.com/milyrock/PR-Reviewer/internal/models"
	"github.com/milyrock/PR-Reviewer/internal/repository"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(repo *repository.Repository) *UserHandler {
	return &UserHandler{service: service.NewUserService(repo)}
}

func (h *UserHandler) SetIsActive(w http.ResponseWriter, r *http.Request) {
	var req models.SetIsActiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, statusBadRequest, errorCodeInvalidRequest, errorMsgInvalidRequestBody)
		return
	}

	user, err := h.service.SetIsActive(req)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user": user,
	})
}

func (h *UserHandler) GetReview(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		writeError(w, statusBadRequest, errorCodeInvalidRequest, errorMsgUserIDRequired)
		return
	}

	prs, err := h.service.GetReview(userID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id":       userID,
		"pull_requests": prs,
	})
}
