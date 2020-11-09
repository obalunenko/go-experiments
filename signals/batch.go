package main

import (
	"errors"
	"time"

	log "github.com/sirupsen/logrus"
)

var closedChan = make(chan struct{})

func init() {
	close(closedChan)
}

var errBatchClosed = errors.New("batch closed")

type MemoryBatch struct {
	items []interface{}
	doFn  BatchDoFn

	maxSize int
	maxWait time.Duration

	flushJobs chan flushRequest
	isRun     bool

	/*notifier channel*/
	insertChan     chan interface{}
	forceFlushChan chan struct{}
	stopChan       chan struct{}
	doneChan       chan struct{}
}

func NewMemoryBatch(flushHandler BatchDoFn, flushMaxSize int, flushMaxWait time.Duration, workerSize int) Batch {
	instance := &MemoryBatch{
		items:          []interface{}{},
		doFn:           flushHandler,
		maxSize:        flushMaxSize,
		maxWait:        flushMaxWait,
		flushJobs:      make(chan flushRequest, workerSize),
		isRun:          false,
		insertChan:     make(chan interface{}, flushMaxSize),
		forceFlushChan: make(chan struct{}),
		stopChan:       make(chan struct{}),
		doneChan:       make(chan struct{}),
	}

	instance.setFlushWorker(workerSize)
	instance.isRun = true
	go instance.run()
	return instance
}

/* Flush Section */
func (mb *MemoryBatch) flush(workerID int, req flushRequest) {
	const maxRetryCount = 3
	var err error

	for i := 0; i < maxRetryCount; i++ {
		if err = mb.doFn(workerID, req.datas); err != nil {
			log.Errorf("failed to execute job: %v", err)
			time.Sleep(500 * time.Millisecond)
		} else {
			break
		}
	}

	if req.doneChan != nil {
		req.doneChan <- struct{}{}
		close(req.doneChan)
	}
}

func (mb *MemoryBatch) setFlushWorker(workerSize int) {
	if workerSize < 1 {
		workerSize = 1
	}
	for id := 1; id <= workerSize; id++ {
		go func(workerID int, flushJobs <-chan flushRequest) {
			log.Infof("worker[%d] started", workerID)

			for j := range flushJobs {
				mb.flush(workerID, j)
			}

			log.Infof("worker[%d] finished", workerID)
		}(id, mb.flushJobs)
	}
}

/* Notifier Section*/

func (mb *MemoryBatch) Insert(data interface{}) error {
	if !mb.isRun {
		return errBatchClosed
	}

	mb.insertChan <- data

	return nil
}

func (mb *MemoryBatch) ForceFlush() error {
	if !mb.isRun {
		return errBatchClosed
	}

	mb.forceFlushChan <- struct{}{}

	return nil
}

func (mb *MemoryBatch) Stop() (chan struct{}, error) {
	log.Infof("stop: initiated")

	if !mb.isRun {
		return closedChan, errBatchClosed
	}

	mb.stopChan <- struct{}{}

	log.Infof("stop: finished")

	return mb.doneChan, nil
}

func (mb *MemoryBatch) run() {
	t := time.NewTicker(mb.maxWait)
	defer t.Stop()

	defer func() {
		log.Infof("run: defer started")
		log.Infof("run: defer not flushed items: %d", len(mb.items))
		if len(mb.items) > 0 {
			mb.flushJobs <- flushRequest{
				datas:    mb.items,
				doneChan: mb.doneChan,
			}

			mb.items = mb.items[:0]
		} else {
			mb.doneChan <- struct{}{}
			close(mb.doneChan)
		}

		close(mb.flushJobs)

		log.Infof("run: defer finished")
	}()

	for mb.isRun {
		select {
		case <-t.C:
			log.Infof("run: ticked")

			if len(mb.items) > 0 {
				mb.flushJobs <- flushRequest{
					datas:    mb.items,
					doneChan: nil,
				}

				mb.items = mb.items[:0]
			}
		case item := <-mb.insertChan:
			log.Infof("run: insert received: [%v]", item)

			mb.items = append(mb.items, item)
			if len(mb.items) >= mb.maxSize {
				mb.flushJobs <- flushRequest{
					datas:    mb.items,
					doneChan: nil,
				}

				mb.items = mb.items[:0]
			}
		case <-mb.forceFlushChan:
			log.Infof("run: force flush received")

			if len(mb.items) > 0 {
				mb.flushJobs <- flushRequest{
					datas:    mb.items,
					doneChan: nil,
				}
				mb.items = mb.items[:0]
			}
		case <-mb.stopChan:
			log.Infof("run: stop received")
			mb.isRun = false

			log.Infof("run: not flushed items: %d", len(mb.items))

			return
		}
	}
}
