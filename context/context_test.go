package main

import (
	"context"
	"errors"
	"testing"
)

func job(ctx context.Context, cancel context.CancelCauseFunc) {
	err := errors.New("test_error")

	cancel(err)
}

func TestCancel(t *testing.T) {
	ctx := context.Background()
	childCtx, cancel := context.WithCancelCause(ctx)

	job(childCtx, cancel)

	t.Logf("parent err: %v", ctx.Err())
	t.Logf("parent cause: %v", context.Cause(ctx))
	t.Logf("child err: %v", childCtx.Err())
	t.Logf("child cause: %v", context.Cause(childCtx))
}