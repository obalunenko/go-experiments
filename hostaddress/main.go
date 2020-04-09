package main

import (
	"os"

	"fmt"
)

func getAddr() (string, error) {
	return os.Hostname()
}

func main() {
	fmt.Println(getAddr())
}
