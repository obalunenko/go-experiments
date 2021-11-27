package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	logger "github.com/sirupsen/logrus"
)

func handler(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()
		key := values.Get("error")
		if key == "" {
			key = "0"
		}

		et, err := strconv.Atoi(key)
		if err != nil {
			err = fmt.Errorf("wrong value of error type [%s]", key)
		}

		if maybeAbortRequest(w, r, err) {
			return
		}

		err = jobEmulation(r.Context(), errorType(et))
		if maybeAbortRequest(w, r, err) {
			return
		}

		fmt.Fprintf(w, "All is ok")
	}
}

type errorType uint

const (
	errorTypeNil errorType = iota
	errorTypeInternal
	errorTypeInvalidValue
	errorTypeDuplicate
	errorTypeNotImplemented
	errorTypeNotFound
	errorTypeTimeout
)

func jobEmulation(ctx context.Context, et errorType) error {
	logger.WithField("error type", et).Info("emulating job with error")

	switch et {
	case errorTypeDuplicate:
		return ErrDuplicate
	case errorTypeInternal:
		return ErrInternal
	case errorTypeNotImplemented:
		return ErrNotImplemented
	case errorTypeTimeout:
		return ErrTimeout
	case errorTypeInvalidValue:
		return ErrInvalidValue
	case errorTypeNotFound:
		return ErrNotFound
	case errorTypeNil:
		return nil
	}

	return nil
}
