package sdkcm

import (
	"net/http"
)

type Response struct {
	StatusCode    int         `json:"status_code,omitempty" form:"status_code"`
	Messages      []string    `json:"messages,omitempty" form:"messages"`
	RootCause     error       `json:"-"`
	InternalError *string     `json:"internal_error"`
	Data          interface{} `json:"data,omitempty" form:"data"`
	Input         interface{} `json:"input,omitempty" form:"input"`
	Paging        *Paging     `json:"paging,omitempty" form:"paging"`
}

func ErrorResponse(statusCode int, input interface{}, rootCause error, messages ...string) *Response {
	var internalError *string
	if rootCause != nil {
		e := rootCause.Error()
		internalError = &e
	}

	if len(messages) <= 0 && rootCause != nil && statusCode > 299 {
		messages = append(messages, rootCause.Error())
	}

	return &Response{
		StatusCode:    statusCode,
		Messages:      messages,
		RootCause:     rootCause,
		Input:         input,
		InternalError: internalError,
	}
}

func SuccessResponse(input, data interface{}, paging *Paging) *Response {
	return &Response{
		StatusCode: http.StatusOK,
		Input:      input,
		Data:       data,
		Paging:     paging,
	}
}

func NewBadRequestResponse(input interface{}, rootCause error, messages ...string) *Response {
	return ErrorResponse(http.StatusBadRequest, input, rootCause, messages...)
}

func NewNotFoundResponse(input interface{}, rootCause error, messages ...string) *Response {
	return ErrorResponse(http.StatusNotFound, input, rootCause, messages...)
}

func NewConflictResponse(input interface{}, rootCause error, messages ...string) *Response {
	return ErrorResponse(http.StatusConflict, input, rootCause, messages...)
}

func NewUnauthorizedResponse(input interface{}, rootCause error, messages ...string) *Response {
	return ErrorResponse(http.StatusUnauthorized, input, rootCause, messages...)
}

func NewInternalServerErrorResponse(input interface{}, rootCause error, messages ...string) *Response {
	return ErrorResponse(http.StatusInternalServerError, input, rootCause, messages...)
}
