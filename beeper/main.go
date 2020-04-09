package main

import (
	beeper "github.com/frozzare/go-beeper"
)

func main() {
	// beep once
	beeper.Beep()

	// beep three times
	beeper.Beep(3)
}
