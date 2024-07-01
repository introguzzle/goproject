package response

import (
	"encoding/json"
	"goproject/src/optional"
	"net/http"
)

type Response interface {
	Serialize() ([]byte, error)
	Status() int
}

type DataResponse struct {
	Data any `json:"data"`
	Code int `json:"status"`
}

type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"status"`
}

func (r *DataResponse) Status() int {
	return r.Code
}

func (r *ErrorResponse) Status() int {
	return r.Code
}

func (r *DataResponse) Serialize() ([]byte, error) {
	return json.Marshal(r)
}

func (r *ErrorResponse) Serialize() ([]byte, error) {
	return json.Marshal(r)
}

func WrapData(
	data any,
	code int,
) *DataResponse {
	return &DataResponse{
		Data: data,
		Code: code,
	}
}

func WrapOptional(
	data optional.Optional[any],
	code int,
) interface{} {
	if data.IsPresent() {
		return WrapData(data.Get(), code)
	} else {
		return WrapErrorString("Not found", code)
	}
}

func Ok() *DataResponse {
	return WrapData("OK", http.StatusOK)
}

func BadRequest() *ErrorResponse {
	return &ErrorResponse{
		Error: "Bad Request",
		Code:  http.StatusBadRequest,
	}
}

func WrapError(
	err error,
	code int,
) *ErrorResponse {
	return &ErrorResponse{
		Error: err.Error(),
		Code:  code,
	}
}

func WrapErrorString(
	err string,
	code int,
) *ErrorResponse {
	return &ErrorResponse{
		Error: err,
		Code:  code,
	}
}

func Internal() *ErrorResponse {
	return &ErrorResponse{
		Error: "Internal server error",
		Code:  http.StatusInternalServerError,
	}
}

func NotFound() *ErrorResponse {
	return &ErrorResponse{
		Error: "Not found",
		Code:  http.StatusNotFound,
	}
}
