package jsonutil

// ErrorResponse for all methods
type ErrorResponse struct {
	Error struct {
		ErrorCode int    `json:"error_code"`
		ErrorMsg  string `json:"error_msg"`
	} `json:"error"`
}

// SuccessfulResponse for all methods returning data
type SuccessfulResponse struct {
	Response interface{} `json:"response"`
}

// NewSuccessfulResponse create and return
func NewSuccessfulResponse(response interface{}) SuccessfulResponse {
	var res SuccessfulResponse

	res.Response = response
	return res
}

// NewError create and return
func NewError(code int, msg string) ErrorResponse {
	var err ErrorResponse

	err.Error.ErrorCode = code
	err.Error.ErrorMsg = msg

	return err
}