package http

import (
	"AvitoTestTask/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"log/slog"
	"net/http"
)

var logger slog.Logger

func init() {
	logger = slog.Logger{}
}

type Handler struct {
	prService   service.PullRequestService
	userService service.UserService
	TeamService service.TeamService
	statistic   service.StatisticService
	router      *chi.Mux
}

func NewHandler(
	pr service.PullRequestService,
	user service.UserService,
	team service.TeamService,
	statistic service.StatisticService,

) *Handler {
	h := &Handler{
		prService:   pr,
		userService: user,
		TeamService: team,
		statistic:   statistic,
		router:      chi.NewRouter(),
	}
	return h
}
func (h *Handler) LoadRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/v1", func(api chi.Router) {

		// Teams
		api.Post("/team/add", h.CreateTeam)
		api.Get("/team/get", h.ListTeamUsers)

		// Users
		api.Post("/users/setIsActive", h.ChangeActivityFlag)
		api.Get("/users/getReview", h.GetUserAssignedPR)

		// PullRequests
		api.Post("/pullRequest/create", h.CreatePR)
		api.Post("/pullRequest/merge", h.MergePR)
		api.Post("/pullRequest/reassign", h.RemoveReviewer)
		
		api.Get("/statistic/user", h.GetUsersStats)
	})

	r.Get("/swagger/*", httpSwagger.WrapHandler)
	return r
}

func (h *Handler) Router() http.Handler {
	r := h.LoadRoutes()
	return r
}

type ErrorResponse struct {
	Error error `json:"error"`
}
