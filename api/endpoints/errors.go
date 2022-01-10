package endpoints

type ErrorResponse struct {
	Error string `json:"error"`
}

var NotFoundResponse = ErrorResponse{
	Error: "not found",
}

func NewErrorResponse(err error) (r *ErrorResponse) {
	return &ErrorResponse{
		Error: err.Error(),
	}
}
