package schemas

type UserSetIsActiveRequest struct {
	ID       string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

type User struct {
	ID       string `json:"user_id"`
	Name     string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

type UserRsponse struct {
	User User `json:"user"`
}

type UserReviewResponse struct {
	UserID string    `json:"user_id"`
	PRs    []PRShort `json:"pull_requests"`
}
