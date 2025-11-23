package dto

type UserStatResponse struct {
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	TotalPRs  int    `json:"total_prs"`
	OpenPRs   int    `json:"open_prs"`
	MergedPRs int    `json:"merged_prs"`
	ClosedPRs int    `json:"closed_prs"`
}

type GetUsersStatsResponse struct {
	Stats []UserStatResponse `json:"stats"`
}
