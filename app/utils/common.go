package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// JSONer if your response implements this interface it'll be called to send as response.
type JSONer interface {
	JSON() interface{}
}

// ResponseJSON respond as a json object.
func ResponseJSON(w http.ResponseWriter, statusCode int, resp interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	if resp == nil {
		return
	}
	if j, ok := resp.(JSONer); ok {
		resp = j.JSON()
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		fmt.Printf("Unable to send response to client: %v", err)
	}
}

// ResponseDTO type
type ResponseDTO struct {
	StatusCode int         `json:"status_code"`
	Data       interface{} `json:"data"`
}

// ResponseSuccess handles the success responses
func ResponseSuccess(
	statusCode int,
	data interface{},
) *ResponseDTO {
	return &ResponseDTO{
		StatusCode: statusCode,
		Data:       data,
	}
}

// ErrorData represents data in ErrorResponse object.
type ErrorData struct {
	ExceptionType    string `json:"exception_type"`
	Message          string `json:"message"`
	DeveloperMessage string `json:"developer_message"`
	MoreInformation  string `json:"more_information"`
}

// ErrorResponse ...
type ErrorResponse struct {
	StatusCode int       `json:"status_code"`
	Data       ErrorData `json:"data"`
}

// NewErrorResponse ...
func NewErrorResponse(statusCode int, errorMsg string) *ErrorResponse {
	e := ErrorResponse{StatusCode: statusCode}
	e.Data.Message = errorMsg
	e.Data.DeveloperMessage = errorMsg
	return &e
}

// NewErrorWithDataResponse instantiates and returns ErrorResponse object.
func NewErrorWithDataResponse(statusCode int, data ErrorData) *ErrorResponse {
	return &ErrorResponse{
		StatusCode: statusCode,
		Data:       data,
	}
}

func (e ErrorResponse) Error() string {
	return e.Data.Message
}

// JSON ...
func (e ErrorResponse) JSON() interface{} {
	return e
}

// ResponseError handles the error responses
func ResponseError(
	statusCode int,
	exType string,
	msg string,
	devMsg string,
	moreInfo string,
) *ErrorResponse {
	return NewErrorWithDataResponse(
		statusCode,
		ErrorData{
			ExceptionType:    exType,
			Message:          msg,
			DeveloperMessage: devMsg,
			MoreInformation:  moreInfo,
		},
	)
}
