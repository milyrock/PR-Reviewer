package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	v1 "github.com/milyrock/PR-Reviewer/internal/handlers/v1"
	"github.com/milyrock/PR-Reviewer/internal/models"
	"github.com/milyrock/PR-Reviewer/internal/repository"
	"github.com/milyrock/PR-Reviewer/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestServer(t *testing.T) (*httptest.Server, func()) {
	ctx := context.Background()
	db, cleanup, err := test.SetupTestDB(ctx)
	require.NoError(t, err)

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

	server := httptest.NewServer(r)

	return server, cleanup
}

func TestHealthEndpoint(t *testing.T) {
	server, cleanup := setupTestServer(t)
	defer cleanup()

	resp, err := http.Get(server.URL + "/health")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestCreateTeamAndPR(t *testing.T) {
	server, cleanup := setupTestServer(t)
	defer cleanup()

	// 1. Создаем команду
	teamReq := models.CreateTeamRequest{
		TeamName: "backend",
		Members: []models.TeamMember{
			{UserID: "u1", Username: "Alice", IsActive: true},
			{UserID: "u2", Username: "Bob", IsActive: true},
			{UserID: "u3", Username: "Charlie", IsActive: true},
		},
	}

	teamBody, _ := json.Marshal(teamReq)
	resp, err := http.Post(server.URL+"/team/add", "application/json", bytes.NewBuffer(teamBody))
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	// 2. Создаем PR
	prReq := models.CreatePRRequest{
		PullRequestID:   "pr-1",
		PullRequestName: "Add feature",
		AuthorID:        "u1",
	}

	prBody, _ := json.Marshal(prReq)
	resp, err = http.Post(server.URL+"/pullRequest/create", "application/json", bytes.NewBuffer(prBody))
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var prResp struct {
		PR models.PullRequest `json:"pr"`
	}
	err = json.NewDecoder(resp.Body).Decode(&prResp)
	require.NoError(t, err)
	resp.Body.Close()

	// 3. Проверяем, что назначены ревьюверы (до 2, исключая автора)
	assert.Equal(t, "OPEN", prResp.PR.Status)
	assert.LessOrEqual(t, len(prResp.PR.AssignedReviewers), 2)
	assert.NotContains(t, prResp.PR.AssignedReviewers, "u1") // Автор не должен быть ревьювером
	assert.Greater(t, len(prResp.PR.AssignedReviewers), 0)   // Должен быть хотя бы один ревьювер
}

func TestMergePRAndReassignRestriction(t *testing.T) {
	server, cleanup := setupTestServer(t)
	defer cleanup()

	// 1. Создаем команду
	teamReq := models.CreateTeamRequest{
		TeamName: "frontend",
		Members: []models.TeamMember{
			{UserID: "u1", Username: "Alice", IsActive: true},
			{UserID: "u2", Username: "Bob", IsActive: true},
		},
	}

	teamBody, _ := json.Marshal(teamReq)
	resp, err := http.Post(server.URL+"/team/add", "application/json", bytes.NewBuffer(teamBody))
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	// 2. Создаем PR
	prReq := models.CreatePRRequest{
		PullRequestID:   "pr-2",
		PullRequestName: "Fix bug",
		AuthorID:        "u1",
	}

	prBody, _ := json.Marshal(prReq)
	resp, err = http.Post(server.URL+"/pullRequest/create", "application/json", bytes.NewBuffer(prBody))
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	// 3. Мержим PR
	mergeReq := models.MergePRRequest{
		PullRequestID: "pr-2",
	}

	mergeBody, _ := json.Marshal(mergeReq)
	resp, err = http.Post(server.URL+"/pullRequest/merge", "application/json", bytes.NewBuffer(mergeBody))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var mergeResp struct {
		PR models.PullRequest `json:"pr"`
	}
	err = json.NewDecoder(resp.Body).Decode(&mergeResp)
	require.NoError(t, err)
	assert.Equal(t, "MERGED", mergeResp.PR.Status)
	resp.Body.Close()

	// 4. Пытаемся переназначить ревьювера - должна быть ошибка
	reassignReq := models.ReassignPRRequest{
		PullRequestID: "pr-2",
		OldUserID:     "u2",
	}

	reassignBody, _ := json.Marshal(reassignReq)
	resp, err = http.Post(server.URL+"/pullRequest/reassign", "application/json", bytes.NewBuffer(reassignBody))
	require.NoError(t, err)
	require.Equal(t, http.StatusConflict, resp.StatusCode)

	var errorResp models.ErrorResponse
	err = json.NewDecoder(resp.Body).Decode(&errorResp)
	require.NoError(t, err)
	assert.Equal(t, "PR_MERGED", errorResp.Error.Code)
	resp.Body.Close()
}

func TestReassignReviewer(t *testing.T) {
	server, cleanup := setupTestServer(t)
	defer cleanup()

	// 1. Создаем команду с 4 пользователями
	teamReq := models.CreateTeamRequest{
		TeamName: "devops",
		Members: []models.TeamMember{
			{UserID: "u1", Username: "Alice", IsActive: true},
			{UserID: "u2", Username: "Bob", IsActive: true},
			{UserID: "u3", Username: "Charlie", IsActive: true},
			{UserID: "u4", Username: "David", IsActive: true},
		},
	}

	teamBody, _ := json.Marshal(teamReq)
	resp, err := http.Post(server.URL+"/team/add", "application/json", bytes.NewBuffer(teamBody))
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	// 2. Создаем PR
	prReq := models.CreatePRRequest{
		PullRequestID:   "pr-3",
		PullRequestName: "Deploy config",
		AuthorID:        "u1",
	}

	prBody, _ := json.Marshal(prReq)
	resp, err = http.Post(server.URL+"/pullRequest/create", "application/json", bytes.NewBuffer(prBody))
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var createResp struct {
		PR models.PullRequest `json:"pr"`
	}
	err = json.NewDecoder(resp.Body).Decode(&createResp)
	require.NoError(t, err)
	originalReviewers := createResp.PR.AssignedReviewers
	require.Greater(t, len(originalReviewers), 0)
	resp.Body.Close()

	// 3. Переназначаем ревьювера
	oldReviewer := originalReviewers[0]
	reassignReq := models.ReassignPRRequest{
		PullRequestID: "pr-3",
		OldUserID:     oldReviewer,
	}

	reassignBody, _ := json.Marshal(reassignReq)
	resp, err = http.Post(server.URL+"/pullRequest/reassign", "application/json", bytes.NewBuffer(reassignBody))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var reassignResp struct {
		PR         models.PullRequest `json:"pr"`
		ReplacedBy string             `json:"replaced_by"`
	}
	err = json.NewDecoder(resp.Body).Decode(&reassignResp)
	require.NoError(t, err)

	// Проверяем, что старый ревьювер заменен
	assert.NotContains(t, reassignResp.PR.AssignedReviewers, oldReviewer)
	assert.Contains(t, reassignResp.PR.AssignedReviewers, reassignResp.ReplacedBy)
	assert.NotEqual(t, oldReviewer, reassignResp.ReplacedBy)
	resp.Body.Close()
}

func TestInactiveUserNotAssigned(t *testing.T) {
	server, cleanup := setupTestServer(t)
	defer cleanup()

	// 1. Создаем команду с активным и неактивным пользователем
	teamReq := models.CreateTeamRequest{
		TeamName: "qa",
		Members: []models.TeamMember{
			{UserID: "u1", Username: "Alice", IsActive: true},
			{UserID: "u2", Username: "Bob", IsActive: false}, // Неактивный
			{UserID: "u3", Username: "Charlie", IsActive: true},
		},
	}

	teamBody, _ := json.Marshal(teamReq)
	resp, err := http.Post(server.URL+"/team/add", "application/json", bytes.NewBuffer(teamBody))
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	// 2. Создаем PR
	prReq := models.CreatePRRequest{
		PullRequestID:   "pr-4",
		PullRequestName: "Test feature",
		AuthorID:        "u1",
	}

	prBody, _ := json.Marshal(prReq)
	resp, err = http.Post(server.URL+"/pullRequest/create", "application/json", bytes.NewBuffer(prBody))
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var prResp struct {
		PR models.PullRequest `json:"pr"`
	}
	err = json.NewDecoder(resp.Body).Decode(&prResp)
	require.NoError(t, err)

	// 3. Проверяем, что неактивный пользователь не назначен
	assert.NotContains(t, prResp.PR.AssignedReviewers, "u2")
	resp.Body.Close()
}

func TestGetReview(t *testing.T) {
	server, cleanup := setupTestServer(t)
	defer cleanup()

	// 1. Создаем команду
	teamReq := models.CreateTeamRequest{
		TeamName: "backend",
		Members: []models.TeamMember{
			{UserID: "u1", Username: "Alice", IsActive: true},
			{UserID: "u2", Username: "Bob", IsActive: true},
		},
	}

	teamBody, _ := json.Marshal(teamReq)
	resp, err := http.Post(server.URL+"/team/add", "application/json", bytes.NewBuffer(teamBody))
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	// 2. Создаем PR
	prReq := models.CreatePRRequest{
		PullRequestID:   "pr-5",
		PullRequestName: "New feature",
		AuthorID:        "u1",
	}

	prBody, _ := json.Marshal(prReq)
	resp, err = http.Post(server.URL+"/pullRequest/create", "application/json", bytes.NewBuffer(prBody))
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	// 3. Получаем список PR для ревьювера
	resp, err = http.Get(server.URL + "/users/getReview?user_id=u2")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var reviewResp struct {
		UserID       string                    `json:"user_id"`
		PullRequests []models.PullRequestShort `json:"pull_requests"`
	}
	err = json.NewDecoder(resp.Body).Decode(&reviewResp)
	require.NoError(t, err)

	assert.Equal(t, "u2", reviewResp.UserID)
	assert.Greater(t, len(reviewResp.PullRequests), 0)
	assert.Equal(t, "pr-5", reviewResp.PullRequests[0].PullRequestID)
	resp.Body.Close()
}

func TestStatisticsEndpoint(t *testing.T) {
	server, cleanup := setupTestServer(t)
	defer cleanup()

	// 1. Создаем команду
	teamReq := models.CreateTeamRequest{
		TeamName: "backend",
		Members: []models.TeamMember{
			{UserID: "u1", Username: "Alice", IsActive: true},
			{UserID: "u2", Username: "Bob", IsActive: true},
		},
	}

	teamBody, _ := json.Marshal(teamReq)
	resp, err := http.Post(server.URL+"/team/add", "application/json", bytes.NewBuffer(teamBody))
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	// 2. Создаем несколько PR
	for i := 1; i <= 3; i++ {
		prReq := models.CreatePRRequest{
			PullRequestID:   fmt.Sprintf("pr-%d", i),
			PullRequestName: fmt.Sprintf("Feature %d", i),
			AuthorID:        "u1",
		}

		prBody, _ := json.Marshal(prReq)
		resp, err = http.Post(server.URL+"/pullRequest/create", "application/json", bytes.NewBuffer(prBody))
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, resp.StatusCode)
		resp.Body.Close()
	}

	// 3. Получаем статистику
	resp, err = http.Get(server.URL + "/statistics")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var statsResp models.StatisticsResponse
	err = json.NewDecoder(resp.Body).Decode(&statsResp)
	require.NoError(t, err)

	// Проверяем, что статистика не пустая
	assert.Greater(t, len(statsResp.UserStats), 0)
	assert.Greater(t, len(statsResp.PRStats), 0)

	// Проверяем, что есть статистика по пользователям
	foundU2 := false
	for _, userStat := range statsResp.UserStats {
		if userStat.UserID == "u2" {
			assert.Greater(t, userStat.ReviewCount, 0)
			foundU2 = true
		}
	}
	assert.True(t, foundU2, "User u2 should be in statistics")

	// Проверяем, что есть статистика по PR
	assert.Equal(t, 3, len(statsResp.PRStats))
	resp.Body.Close()
}

func TestMergeIdempotency(t *testing.T) {
	server, cleanup := setupTestServer(t)
	defer cleanup()

	// 1. Создаем команду
	teamReq := models.CreateTeamRequest{
		TeamName: "backend",
		Members: []models.TeamMember{
			{UserID: "u1", Username: "Alice", IsActive: true},
			{UserID: "u2", Username: "Bob", IsActive: true},
		},
	}

	teamBody, _ := json.Marshal(teamReq)
	resp, err := http.Post(server.URL+"/team/add", "application/json", bytes.NewBuffer(teamBody))
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	// 2. Создаем PR
	prReq := models.CreatePRRequest{
		PullRequestID:   "pr-6",
		PullRequestName: "Idempotent merge",
		AuthorID:        "u1",
	}

	prBody, _ := json.Marshal(prReq)
	resp, err = http.Post(server.URL+"/pullRequest/create", "application/json", bytes.NewBuffer(prBody))
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	// 3. Мержим PR первый раз
	mergeReq := models.MergePRRequest{
		PullRequestID: "pr-6",
	}

	mergeBody, _ := json.Marshal(mergeReq)
	resp, err = http.Post(server.URL+"/pullRequest/merge", "application/json", bytes.NewBuffer(mergeBody))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()

	// 4. Мержим PR второй раз - должно быть идемпотентно
	resp, err = http.Post(server.URL+"/pullRequest/merge", "application/json", bytes.NewBuffer(mergeBody))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode) // Не должно быть ошибки

	var mergeResp struct {
		PR models.PullRequest `json:"pr"`
	}
	err = json.NewDecoder(resp.Body).Decode(&mergeResp)
	require.NoError(t, err)
	assert.Equal(t, "MERGED", mergeResp.PR.Status)
	resp.Body.Close()
}
