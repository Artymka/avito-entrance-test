package response

const (
	CodeServer     = "SERVER_ERROR"
	CodeBadRequest = "BAD_REQUEST"

	CodeTeamExists  = "TEAM_EXISTS"
	CodePRExists    = "PR_EXISTS"
	CodePRMerged    = "PR_MERGED"
	CodeNotAssigned = "NOT_ASSIGNED"
	CodeNoCandidate = "NO_CANDIDATE"
	CodeNotFound    = "NOT_FOUND"
)

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Err Error `json:"error"`
}

func Err(code string, message string) ErrorResponse {
	return ErrorResponse{
		Err: Error{
			Code:    code,
			Message: message,
		},
	}
}
