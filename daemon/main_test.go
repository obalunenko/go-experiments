package main

import (
	"testing"

	"go.uber.org/goleak"
)

func Test_run(t *testing.T) {
	defer goleak.VerifyNone(t)

	run()
}
