package main

import (
	"bufio"
	"context"
	"os"
	"os/signal"
	"time"

	log "github.com/sirupsen/logrus"
)

func main() {
	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt, os.Kill)

	b := NewMemoryBatch(job, 10, 30*time.Second, 3)

	ctx, cancel := context.WithCancel(context.Background())

	go handleInput(ctx, b)

	sig := <-c

	log.Infof("received signal: [%s]", sig.String())
	cancel()
	done, err := b.Stop()

	if err != nil {
		log.Errorf("failed to stop batch: %v", err)
	}

	<-done

	log.Info("Exit")
}

func job(workerID int, data []interface{}) error {
	log.Infof("[worker:%d] - in job", workerID)
	for _, d := range data {
		log.Infof("[worker:%d] - %v \n", workerID, d)
	}

	log.Infof("[worker:%d] - sleep to imitate lag)", workerID)

	time.Sleep(time.Second)

	log.Infof("[worker:%d] - job finished", workerID)

	return nil
}

func handleInput(ctx context.Context, b Batch) {
	input := make(chan interface{}, 1)

	go func(ctx context.Context, out chan<- interface{}) {
		reader := bufio.NewReader(os.Stdin)
		scanner := bufio.NewScanner(reader)

		for scanner.Scan() {
			line := scanner.Text()
			if err := scanner.Err(); err != nil {
				log.Errorf("handleInput: scan error: %v \n", err)
			}

			select {
			case <-ctx.Done():
				close(out)
				return
			case out <- line:
				log.Infof("handleInput: line sent to input: [%s] \n", line)
			}
		}

	}(ctx, input)

	for {
		select {
		case <-ctx.Done():
			log.Infof("handleInput: done context")
			return
		case in := <-input:
			log.Infof("handleInput: from input [%s]", in)
			if err := b.Insert(in); err != nil {
				log.Errorf("failed to insert to batch: %v", err)
			}
		}
	}
}
