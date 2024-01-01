package handler

import "net/http"

type ResultResponse[T any] struct {
	Status     string `json:"status"`
	StatusCode int    `json:"status_code"`
	Result     T      `json:"result"`
} // @name ResultResponse

type MessageResponse struct {
	Status     string `json:"status"`
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
} // @name MessageResponse

type ErrorResponse struct {
	Status     string `json:"status"`
	StatusCode int    `json:"status_code"`
	Error      string `json:"error"`
} // @name ErrorResponse

func makeHttpResponse[T any](code int, result T) (int, ResultResponse[T]) {
	response := ResultResponse[T]{
		Status:     http.StatusText(code),
		StatusCode: code,
		Result:     result,
	}

	return code, response
}

func makeHttpMessageResponse(code int, message string) (int, MessageResponse) {
	response := MessageResponse{
		Status:     http.StatusText(code),
		StatusCode: code,
		Message:    message,
	}

	return code, response
}

func makeHttpErrorResponse(code int, err string) (int, ErrorResponse) {
	response := ErrorResponse{
		Status:     http.StatusText(code),
		StatusCode: code,
		Error:      err,
	}

	return code, response
}
