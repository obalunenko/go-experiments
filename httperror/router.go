package main

import (
	"context"
	"net/http"
)

func NewRouter(ctx context.Context) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/job", handler(ctx))

	return mux
}
