package main

import (
	"fmt"
	"os"
)

func getAddr() (string, error) {
	return os.Hostname()
}

func main() {
	fmt.Println(getAddr())
}
