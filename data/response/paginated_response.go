package response

type PaginatedResponse struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
	Limit  int         `json:"limit"`
	Offset int         `json:"offset"`
	Msg    string      `json:"msg"`
}