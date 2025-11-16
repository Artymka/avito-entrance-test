package schemas

type TeamAddRequest struct {
	Name    string   `json:"team_name"`
	Members []Member `json:"members"`
}

type Member struct {
	ID       string `json:"user_id"`
	Name     string `json:"username"`
	IsActive bool   `json:"is_active"`
}
