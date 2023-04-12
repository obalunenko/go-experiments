package main

import (
	"context"
	"errors"
	"testing"
	"time"
)

func jobWithCause(ctx context.Context, cancel context.CancelCauseFunc) {
	err := errors.New("test_error")

	cancel(err)
}

func TestCancelWithCause(t *testing.T) {
	ctx := context.Background()
	childCtx, cancel := context.WithCancelCause(ctx)

	defer cancel(nil)

	jobWithCause(childCtx, cancel)

	t.Logf("parent err: %v", ctx.Err())
	t.Logf("parent cause: %v", context.Cause(ctx))
	t.Logf("child err: %v", childCtx.Err())
	t.Logf("child cause: %v", context.Cause(childCtx))
}

func TestCancel(t *testing.T) {
	ctx := context.Background()
	childCtx, cancel := context.WithCancel(ctx)

	childCtx, cancel = context.WithDeadline(childCtx, time.Date(2022, 12, 12, 0, 0, 0, 0, time.UTC))

	cancel()

	time.Sleep(time.Second * 2)

	t.Logf("parent err: %v", ctx.Err())
	t.Logf("parent cause: %v", context.Cause(ctx))
	t.Logf("child err: %v", childCtx.Err())
	t.Logf("child cause: %v", context.Cause(childCtx))
}
