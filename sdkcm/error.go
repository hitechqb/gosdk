package sdkcm

import "net/http"

type SDKError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	RootCause  error  `json:"root_cause"`
}

func (s SDKError) Error() string {
	return s.Message
}

func (s SDKError) String() string {
	return s.Message
}

func NewSDKError(statusCode int, message string, rootCause error) SDKError {
	return SDKError{StatusCode: statusCode, Message: message, RootCause: rootCause}
}

var (
	ErrInvalidRequest = NewSDKError(http.StatusBadRequest, "Thông tin không hợp lệ", nil)
	ErrDataNotFound   = NewSDKError(http.StatusNotFound, "Dữ liệu không tìm thấy", nil)
	ErrDataConflict   = NewSDKError(http.StatusConflict, "Dữ liệu đã tồn tại trong hệ thống", nil)
	ErrDB             = func(rootCause error) error {
		return NewSDKError(http.StatusUnprocessableEntity, "Không thể truy vấn dữ liệu này", rootCause)
	}
	ErrBadRequest = func(message string, rootCause error) error {
		return NewSDKError(http.StatusBadRequest, message, rootCause)
	}
)
