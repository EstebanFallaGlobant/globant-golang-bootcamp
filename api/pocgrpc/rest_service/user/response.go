package user

type errorResponse struct {
	Message string `json:"msg"`
	Code    int    `json:"status"`
}
