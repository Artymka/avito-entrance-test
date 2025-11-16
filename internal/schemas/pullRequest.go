package schemas

type PRCreateRequest struct {
	ID       string `json:"pull_request_id"`
	Name     string `json:"pull_request_name"`
	AuthorID string `json:"author_id"`
}

type PRResponse struct {
	ID               string   `json:"pull_request_id"`
	Name             string   `json:"pull_request_name"`
	AuthorID         string   `json:"author_id"`
	Status           string   `json:"status"`
	AssignedReviwers []string `json:"assigned_reviewers"`
}

type PRCreateResponse struct {
	PR PRResponse `json:"pr"`
}

type PRMergeRequest struct {
	ID string `json:"pull_request_id"`
}

type PRMergeResponse struct {
	PR PRMerge `json:"pr"`
}

type PRMerge struct {
	ID               string   `json:"pull_request_id"`
	Name             string   `json:"pull_request_name"`
	AuthorID         string   `json:"author_id"`
	Status           string   `json:"status"`
	AssignedReviwers []string `json:"assigned_reviewers"`
	MergedAt         string   `json:"mergedAt"`
}
