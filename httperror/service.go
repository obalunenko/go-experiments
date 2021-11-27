package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	logger "github.com/sirupsen/logrus"
)

type Service struct {
	appServer *Server
	wg        *sync.WaitGroup
	stopChan  chan os.Signal
	ctx       struct {
		val        context.Context
		cancelFunc context.CancelFunc
	}
}

func (s *Service) Run() chan struct{} {
	s.wg.Add(1)

	go s.appServer.Run()

	doneChan := make(chan struct{})

	go func() {
	loop:
		for {
			select {
			case sig := <-s.stopChan:
				logger.WithField("signal", sig.String()).Warn("Signal received")

				break loop
			case err := <-s.appServer.Errors():
				if err != nil {
					logger.WithError(err).Error("server error")

					break loop
				}
			}
		}

		s.ctx.cancelFunc()

		s.wg.Wait()

		close(doneChan)
	}()

	return doneChan
}

type ServiceParams struct {
	AppPort string
}

func NewService(ctx context.Context, p ServiceParams) (*Service, error) {
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(ctx)

	logWr := logger.StandardLogger().Writer()

	router := NewRouter(ctx)

	wg.Add(1)

	srv := NewServer(ctx, &wg, "server-example", "8080", logWr, router,
		func(wg *sync.WaitGroup, s *http.Server) {
			defer wg.Done()

			s.ErrorLog.Println("Disable keep-alive")

			s.SetKeepAlivesEnabled(false)

			if err := logWr.Close(); err != nil {
				s.ErrorLog.Printf("failed to close log writer: %v", err)
			}
		},
	)

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT)

	return &Service{
		appServer: srv,
		wg:        &wg,
		stopChan:  stopChan,
		ctx: struct {
			val        context.Context
			cancelFunc context.CancelFunc
		}{
			val:        ctx,
			cancelFunc: cancel,
		},
	}, nil
}
