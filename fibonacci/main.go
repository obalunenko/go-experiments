package main

import (
	"fmt"
	"sync"
	"time"
)

func fibonacci(wg *sync.WaitGroup, jobs chan int, res chan int64) {
	var x, y int64
	x, y = 0, 1
	for range jobs {
		select {
		case res <- x:
			x, y = y, x+y
		}
	}
	close(res)
	wg.Done()
}

func main() {
	jobs := make(chan int)
	res := make(chan int64)
	var wg sync.WaitGroup

	wg.Add(3)
	start := time.Now()
	go job(&wg, jobs)
	go fibonacci(&wg, jobs, res)
	go results(&wg, res)
	wg.Wait()
	elapsed := time.Since(start)
	fmt.Println("quit")
	fmt.Printf("took: %s", elapsed)
}

func job(wg *sync.WaitGroup, jobs chan int) {
	for i := 0; i < 100000; i++ {
		jobs <- i
	}
	close(jobs)
	wg.Done()

}

func results(wg *sync.WaitGroup, res chan int64) {
	for r := range res {
		fmt.Println(r)
	}
	wg.Done()
}
