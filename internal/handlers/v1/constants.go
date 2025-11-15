package v1

import "net/http"

const (
	statusBadRequest    = http.StatusBadRequest          // 400
	statusNotFound      = http.StatusNotFound            // 404
	statusInternalError = http.StatusInternalServerError // 500
	statusConflict      = http.StatusConflict            // 409
	statusCreated       = http.StatusCreated             // 201
)

const (
	errorCodeInvalidRequest = "INVALID_REQUEST"
	errorCodeNotFound       = "NOT_FOUND"
	errorCodeInternalError  = "INTERNAL_ERROR"
	errorCodeTeamExists     = "TEAM_EXISTS"
	errorCodePRExists       = "PR_EXISTS"
	errorCodePRMerged       = "PR_MERGED"
	errorCodeNotAssigned    = "NOT_ASSIGNED"
	errorCodeNoCandidate    = "NO_CANDIDATE"
)

const (
	errorMsgInvalidRequestBody   = "invalid request body"
	errorMsgResourceNotFound     = "resource not found"
	errorMsgTeamNameExists       = "team_name already exists"
	errorMsgPRIDExists           = "PR id already exists"
	errorMsgCannotReassignMerged = "cannot reassign on merged PR"
	errorMsgReviewerNotAssigned  = "reviewer is not assigned to this PR"
	errorMsgNoCandidate          = "no active replacement candidate in team"
	errorMsgTeamNameRequired     = "team_name parameter is required"
	errorMsgUserIDRequired       = "user_id parameter is required"
)
