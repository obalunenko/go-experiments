package main

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"
)

func main() {
	lock := make(chan struct{}, 1)
	msg := make(chan string)

	workersNum := runtime.NumCPU()

	var wg sync.WaitGroup

	wg.Add(workersNum)

	for i := 1; i <= workersNum; i++ {
		go work(&wg, i, lock, msg)
	}

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})

	go printMsg(ctx, msg, done)

	wg.Wait()

	cancel()

	<-done

	log.Println("Exit.")
}

func work(wg *sync.WaitGroup, id int, lock chan struct{}, msg chan<- string) {
	defer wg.Done()

	log.Printf("[%d]: want to lock \n", id)
	lock <- struct{}{}

	log.Printf("[%d]: has lock \n", id)

	time.Sleep(time.Millisecond * 150)

	msg <- fmt.Sprintf("Msg from worker [%d]", id)

	<-lock
	log.Printf("[%d]: unlocked \n", id)
}

func printMsg(ctx context.Context, msg <-chan string, done chan<- struct{}) {
	defer func() {
		done <- struct{}{}
		close(done)
	}()

	for {
		select {
		case <-ctx.Done():
			log.Println("context canceled")
			return
		case m := <-msg:
			fmt.Println(m)
		}
	}
}
