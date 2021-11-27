package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	logger "github.com/sirupsen/logrus"
)

var (
	ErrInternal       = errors.New("internal server error")
	ErrInvalidValue   = errors.New("invalid value")
	ErrDuplicate      = errors.New("duplicated")
	ErrNotImplemented = errors.New("not implemented")
	ErrNotFound       = errors.New("not found")
	ErrTimeout        = errors.New("timeout")
)

// NewHTTPError creates HTTPError with passed status and use error as a message.
func NewHTTPError(status int, err error) HTTPError {
	return HTTPError{
		Code:    status,
		Message: err.Error(),
	}
}

// HTTPError represents http error response message.
type HTTPError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
}

// JSONError is like the http.Error, but allow to pass json to body.
func JSONError(ctx context.Context, w http.ResponseWriter, err interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(err); err != nil {
		logger.WithError(err).Error("failed to encode error message")
	}
}

func maybeAbortRequest(w http.ResponseWriter, r *http.Request, err error) bool {
	if err == nil {
		return false
	}

	ctx := r.Context()

	logger.WithError(err).Error(r.URL.RequestURI())

	var (
		status int
		msg    string
	)

	switch {
	case errors.Is(err, ErrInternal):
		status = http.StatusInternalServerError
		msg = "Some internal server issue"

	case errors.Is(err, ErrInvalidValue):
		status = http.StatusBadRequest
		msg = "Invalid value passed"

	case errors.Is(err, ErrDuplicate):
		status = http.StatusConflict
		msg = "Already exist"

	case errors.Is(err, ErrNotImplemented):
		status = http.StatusNotImplemented
		msg = "Not implemented yet"

	case errors.Is(err, ErrNotFound):
		status = http.StatusNotFound
		msg = "Requested resource not found"

	case errors.Is(err, ErrTimeout):
		status = http.StatusRequestTimeout
		msg = "Request was canceled due timeout reached"

	default:
		status = http.StatusInternalServerError
		msg = "Some internal server issue"
	}

	JSONError(ctx, w, NewHTTPError(status, fmt.Errorf("%s: %w", msg, err)), status)

	return true
}
