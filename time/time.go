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

	t1 := time.Time{}
	t2 := time.Time{}

	if t1.Before(t2) {
		fmt.Println("t1 is before t2")
	} else {
		fmt.Println("t1 is not before t2")
	}
}
