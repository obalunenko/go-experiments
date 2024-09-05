package main

import (
	"github.com/gen2brain/beeep"
)

func main() {
	err := beeep.Notify("Test notification",
		"This is a test of toast notifications",
		"assets/information.png")
	if err != nil {
		return
	}
}
