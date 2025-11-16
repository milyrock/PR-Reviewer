package v1

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/milyrock/PR-Reviewer/internal/models"
	"github.com/milyrock/PR-Reviewer/internal/service"
)

func handleServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrPRExists):
		writeError(w, statusConflict, errorCodePRExists, errorMsgPRIDExists)
	case errors.Is(err, service.ErrPRNotFound), errors.Is(err, service.ErrUserNotFound), errors.Is(err, service.ErrTeamNotFound):
		writeError(w, statusNotFound, errorCodeNotFound, errorMsgResourceNotFound)
	case errors.Is(err, service.ErrPRMerged):
		writeError(w, statusConflict, errorCodePRMerged, errorMsgCannotReassignMerged)
	case errors.Is(err, service.ErrReviewerNotAssigned):
		writeError(w, statusConflict, errorCodeNotAssigned, errorMsgReviewerNotAssigned)
	case errors.Is(err, service.ErrNoCandidate):
		writeError(w, statusConflict, errorCodeNoCandidate, errorMsgNoCandidate)
	case errors.Is(err, service.ErrTeamExists):
		writeError(w, statusBadRequest, errorCodeTeamExists, errorMsgTeamNameExists)
	default:
		writeError(w, statusInternalError, errorCodeInternalError, err.Error())
	}
}

func writeError(w http.ResponseWriter, statusCode int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(models.ErrorResponse{
		Error: struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		}{
			Code:    code,
			Message: message,
		},
	}); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}
