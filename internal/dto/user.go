package dto

type User struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

// ---------- Request ----------
type SetUserIsActiveRequest struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

type GetUserReviewsRequest struct {
	UserID string `json:"user_id"`
}

// ---------- Responce ----------
type SetUserIsActiveResponse struct {
	User User `json:"user"`
}

type GetUserReviewsResponse struct {
	UserID       string             `json:"user_id"`
	PullRequests []PullRequestShort `json:"pull_requests"`
}
