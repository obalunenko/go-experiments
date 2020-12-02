package main

import (
	"fmt"
	"testing"
)

func TestEnums(t *testing.T) {
	fmt.Println(Platform(0).String())
	fmt.Println(Platform(3).String())
}
