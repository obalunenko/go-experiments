package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type daemon struct {
	dataChan chan string
	stopChan chan struct{}
	wg       *sync.WaitGroup
}

func (d *daemon) start(ctx context.Context) {
	fmt.Println("		start daemon")
	defer func() {
		fmt.Println("d.start defer")
		d.wg.Done()
	}()

	for {
		select {
		case <-d.stopChan:
			fmt.Println("		stop daemon signal received")
			return

		case data, open := <-d.dataChan:
			if !open {
				fmt.Println("		data chan was closed")
				return
			}
			fmt.Printf("		daemon received new data %s \n", data)

		case <-ctx.Done():
			fmt.Println("context canceled - stop the daemon")
			return
		}

		fmt.Println("new loop daemon")
	}
}

func (d *daemon) stop() {
	fmt.Println("		d.Stop")
	select {
	case d.stopChan <- struct{}{}:
		fmt.Printf("send d.stopChan \n")
	default:
		fmt.Printf("d.stopChan is full \n")
	}
}

func newDaemon(wg *sync.WaitGroup) *daemon {
	return &daemon{
		dataChan: make(chan string),
		stopChan: make(chan struct{}, 1),
		wg:       wg,
	}

}

// Service ...
type Service struct {
	wg        *sync.WaitGroup
	daemonsWG *sync.WaitGroup
	d         *daemon
	ctx       struct {
		ctx    context.Context
		cancel context.CancelFunc
	}
}

// Start starts service, underlying daemon and start processing of data
func (svc *Service) Start() {
	go svc.StartDaemon()
	fmt.Println("	service started")

	go svc.processData()
}

// Stop stops the whole service
func (svc *Service) Stop() {
	defer svc.wg.Done()

	defer func() {
		close(svc.d.dataChan)
	}()

	fmt.Println("	svc.Stop")
	svc.ctx.cancel()

	// wait for daemons stop
	svc.daemonsWG.Wait()
	fmt.Println("all daemons stopped")
}

// StartDaemon starts the daemon for processing the data
func (svc *Service) StartDaemon() {
	svc.daemonsWG.Add(1)
	go svc.d.start(svc.ctx.ctx)
}

// StopDaemon stops the daemon
func (svc *Service) StopDaemon() {

	svc.d.stop()
}

func (svc *Service) processData() {
	defer func() {
		fmt.Println("svc.ProcessData defer")
	}()

	t := time.NewTicker(time.Second)
	for {
		select {
		case tick, open := <-t.C:
			if !open {
				fmt.Println("	ticker was closed")
				return
			}

			fmt.Println("	ticked")
			data := grabData(tick)

			select {
			case svc.d.dataChan <- data:
				fmt.Println("	send data to daemon")
			default:
				fmt.Println("	dropping the data, daemon stopped")
			}

		case <-svc.ctx.ctx.Done():
			fmt.Println("	service received stop signal")
			t.Stop()
			return
		}

		fmt.Println("new loop process")
	}
}

// NewService creates new instance of Service
func NewService(wg *sync.WaitGroup) *Service {
	var dWg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	return &Service{
		d:         newDaemon(&dWg),
		wg:        wg,
		daemonsWG: &dWg,
		ctx: struct {
			ctx    context.Context
			cancel context.CancelFunc
		}{
			ctx:    ctx,
			cancel: cancel,
		},
	}

}

func grabData(t time.Time) string {
	return t.String()
}

func run() {
	var wg sync.WaitGroup
	wg.Add(1)

	s := NewService(&wg)
	sleep := 5 * time.Second
	fmt.Println("starting service")

	s.Start()
	time.Sleep(sleep)

	s.StopDaemon()
	time.Sleep(sleep)

	s.StartDaemon()
	time.Sleep(sleep)

	s.Stop()

	wg.Wait()
	fmt.Println("service stopped")
	time.Sleep(sleep)
}

func main() {
	run()
}
