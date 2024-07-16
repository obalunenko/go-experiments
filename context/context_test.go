package main

import (
	"context"
	"errors"
	"fmt"
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

	jobWithCause(childCtx, cancel)

	cancel(nil)

	t.Logf("parent err: %v", ctx.Err())
	t.Logf("parent cause: %v", context.Cause(ctx))
	t.Logf("child err: %v", childCtx.Err())
	t.Logf("child cause: %v", context.Cause(childCtx))
}

func TestCancelCause(t *testing.T) {
	ctx := context.Background()
	childCtx, cancel := context.WithCancel(ctx)

	deadLine := time.Date(2022, 12, 12, 0, 0, 0, 0, time.UTC)
	childCtx, cancel = context.WithDeadlineCause(childCtx, deadLine, errors.New("test_deadline"))

	cancel()

	time.Sleep(time.Second * 2)

	t.Logf("parent err: %v", ctx.Err())
	t.Logf("parent cause: %v", context.Cause(ctx))
	t.Logf("child err: %v", childCtx.Err())
	t.Logf("child cause: %v", context.Cause(childCtx))
}

func TestCancel(t *testing.T) {
	ctx := context.Background()
	childCtx, cancel := context.WithCancel(ctx)

	defer cancel()

	tCtx, tCancel := context.WithTimeout(childCtx, time.Second)
	defer tCancel()

	jobWithTime(childCtx, time.Second*2)

	select {
	case <-tCtx.Done():
	}

	t.Logf("parent err: %v", ctx.Err())
	t.Logf("parent cause: %v", context.Cause(ctx))
	t.Logf("child err: %v", childCtx.Err())
	t.Logf("child cause: %v", context.Cause(childCtx))
	t.Logf("tCtx err: %v", tCtx.Err())
	t.Logf("tCtx cause: %v", context.Cause(tCtx))
}

func jobWithTime(ctx context.Context, d time.Duration) {
	if ctx.Err() != nil {
		fmt.Println("jobWithTime: ", ctx.Err())
		return
	}

	time.Sleep(d)
}
