package router

import (
	"net/http"
)

// ErrorPayload represents optional payload data format passed to HTTPError
type ErrorPayload map[string]interface{}

// HTTPError Generic HTTP Error
type HTTPError struct {
	Message string
	Code    int
	Payload ErrorPayload
}

// GetResponseData returns serialized JSON which is sent to client in case of error
func (httpError *HTTPError) GetResponseData() []byte {
	res := map[string]interface{}{
		"statusCode": httpError.Code,
		"error":      httpError.Message,
	}
	if httpError.Payload != nil {
		res["payload"] = httpError.Payload
	}
	jsonByteArr, _ := jsonSerialize(res)
	return jsonByteArr
}

func createError(message string, status int, payload ErrorPayload) *HTTPError {
	return &HTTPError{
		Message: message,
		Code:    status,
		Payload: payload,
	}
}

func errorMessage(message string, msg []string) string {
	if len(msg) > 0 {
		return msg[0]
	}
	return message
}

// BadRequestError HTTP 400
func BadRequestError(payload ErrorPayload, msg ...string) *HTTPError {
	message := errorMessage("Bad Request", msg)
	return createError(message, http.StatusBadRequest, payload)
}

// UnauthorizedError HTTP 401
func UnauthorizedError(payload ErrorPayload, msg ...string) *HTTPError {
	message := errorMessage("Unauthorized", msg)
	return createError(message, http.StatusUnauthorized, payload)
}

// NotFoundError HTTP 404
func NotFoundError(msg ...string) *HTTPError {
	message := errorMessage("Not Found", msg)
	return createError(message, http.StatusNotFound, nil)
}

// MethodNotAllowedError HTTP 405
func MethodNotAllowedError() *HTTPError {
	return createError("Method Not Allowed", http.StatusMethodNotAllowed, nil)
}

// InternalError HTTP 500
func InternalError(msg ...string) *HTTPError {
	message := errorMessage("Internal Server Error", msg)
	return createError(message, http.StatusInternalServerError, nil)
}
