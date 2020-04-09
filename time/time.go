package main

import (
	"fmt"
	"time"
)

func main() {

	nowInLoc := time.Now().In(time.FixedZone("UTC+12", 12*60*60))
	fmt.Println(nowInLoc)
	fmt.Println(nowInLoc.Format(time.RFC3339))

	fmt.Println()

	now := time.Now()
	fmt.Println(now)
	fmt.Println(now.Format(time.RFC3339))

}
