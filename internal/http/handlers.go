package http

import (
	"AvitoTestTask/internal/dto"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// CreateTeam godoc
// @Summary Создать команду с участниками (создаёт/обновляет пользователей)
// @Tags Teams
// @Accept json
// @Produce json
// @Param request body dto.CreateTeamRequest true "Данные команды"
// @Success 201 {object} dto.CreateTeamResponse
// @Failure 400 {object} ErrorResponse
// @Router /team/add [post]
func (h *Handler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var req dto.CreateTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			http.Error(w, `{"error": "request timeout"}`, http.StatusRequestTimeout)
		default:
			http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusNotFound)

		}
		return
	}
	resp, err := h.TeamService.CreateTeam(ctx, req.TeamName, req.Members)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

// ListTeamUsers godoc
// @Summary Получить команду с участниками
// @Tags Teams
// @Accept json
// @Produce json
// @Param team_name query string true "Название команды"
// @Success 200 {object} dto.GetTeamResponse
// @Failure 404 {object} ErrorResponse
// @Router /team/get [get]
func (h *Handler) ListTeamUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	teamName := r.URL.Query().Get("team_name")
	if teamName == "" {
		http.Error(w, `{"error": "team_name is required"}`, http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	resp, err := h.TeamService.ListTeamUsers(ctx, teamName)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			http.Error(w, `{"error": "request timeout"}`, http.StatusRequestTimeout)
		default:
			http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusNotFound)
		}
		return
	}
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(resp); err != nil {
		logger.Error("failed to encode response", "error", err)
	}
}

// ChangeActivityFlag godoc
// @Summary Установить флаг активности пользователя
// @Tags Users
// @Accept json
// @Produce json
// @Param request body dto.SetUserIsActiveRequest true "ID юзера и флаг активности"
// @Success 200 {object} dto.SetUserIsActiveResponse
// @Failure 404 {object} ErrorResponse
// @Router /users/setIsActive [post]
func (h *Handler) ChangeActivityFlag(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var req dto.SetUserIsActiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	user, err := h.userService.ChangeActivityFlag(ctx, req.UserID, req.IsActive)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			http.Error(w, `{"error": "request timeout"}`, http.StatusRequestTimeout)
		default:
			http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusNotFound)

		}
	}
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(user); err != nil {
		logger.Error("failed to encode response", "error", err)
	}
}

// GetUserAssignedPR godoc
// @Summary Получить PR'ы, где пользователь назначен ревьювером
// @Tags Users
// @Accept json
// @Produce json
// @Param user_id query string true "ID пользователя"
// @Success 200 {object} dto.GetUserReviewsResponse
// @Failure 404 {object} ErrorResponse
// @Router /users/getReview [get]
func (h *Handler) GetUserAssignedPR(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	UserId := r.URL.Query().Get("user_id")
	if UserId == "" {
		http.Error(w, `{"error": "user_id is required"}`, http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	response, err := h.userService.GetUserAssignedPR(ctx, UserId)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			http.Error(w, `{"error": "request timeout"}`, http.StatusRequestTimeout)
		default:
			http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusNotFound)
		}
	}
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(response); err != nil {
		logger.Error("failed to encode response", "error", err)
	}
}

// CreatePR godoc
// @Summary Создать PR и автоматически назначить до 2 ревьюверов из команды автора
// @Tags PullRequests
// @Accept json
// @Produce json
// @Param request body dto.CreatePullRequestRequest true "Данные PR"
// @Success 201 {object} dto.PullRequest
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /pullRequest/create [post]
func (h *Handler) CreatePR(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var req dto.CreatePullRequestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	pr, err := h.prService.CreatePR(ctx, req.PullRequestID, req.AuthorID, req.PullRequestName)

	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			http.Error(w, `{"error": "request timeout"}`, http.StatusRequestTimeout)
		default:
			http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusCreated)

	if err = json.NewEncoder(w).Encode(pr); err != nil {
		logger.Error("failed to encode response", "error", err)
	}
}

// MergePR godoc
// @Summary Пометить PR как MERGED (идемпотентно)
// @Tags PullRequests
// @Accept json
// @Produce json
// @Param request body dto.MergePullRequestRequest true "ID PR"
// @Success 200 {object} dto.PullRequest
// @Failure 404 {object} ErrorResponse
// @Router /pullRequest/merge [post]
func (h *Handler) MergePR(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var req dto.MergePullRequestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	pr, err := h.prService.MergePR(ctx, req.PullRequestID)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			http.Error(w, `{"error": "request timeout"}`, http.StatusRequestTimeout)
		default:
			http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)

		}
		return
	}
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(pr)
	if err != nil {
		logger.Error("failed to encode response", "error", err)
	}

}

// RemoveReviewer godoc
// @Summary Переназначить ревьювера на другого из команды
// @Tags PullRequests
// @Accept json
// @Produce json
// @Param request body dto.ReassignReviewerRequest true "Данные переназначения"
// @Success 200 {object} dto.ReassignReviewerResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /pullRequest/reassign [post]
func (h *Handler) RemoveReviewer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var req dto.ReassignReviewerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	pr, newReviewer, err := h.prService.ReassignReviewer(ctx, req.PullRequestID, req.OldUserID)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			http.Error(w, `{"error": "request timeout"}`, http.StatusRequestTimeout)
		default:
			http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusOK)

	response := struct {
		PR         interface{} `json:"pr"`
		ReplacedBy string      `json:"replaced_by"`
	}{
		PR:         pr,
		ReplacedBy: newReviewer,
	}
	if err = json.NewEncoder(w).Encode(response); err != nil {
		logger.Error("failed to encode response", "error", err)
	}
}

// GetUsersStats godoc
// @Summary Получение статистики по пользователям
// @Description Возвращает статистику по всем пользователям: количество review, open PR и merged PR
// @Tags Statistic
// @Accept json
// @Produce json
// @Success 200 {array} dto.GetUsersStatsResponse
// @Failure 500 {object} ErrorResponse
// @Router /statistic/user [get]
func (h *Handler) GetUsersStats(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stats, err := h.statistic.GetUsersStats(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(stats)
	if err != nil {
		logger.Error("failed to encode response", "error", err)
	}

}
