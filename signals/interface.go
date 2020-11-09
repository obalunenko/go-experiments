package main

type Batch interface {
	flush(workerID int, req flushRequest)
	setFlushWorker(workerSize int)

	Insert(data interface{}) (err error)
	ForceFlush() (err error)
	Stop() (doneChan chan struct{}, err error)
}

type BatchDoFn func(workerID int, datas []interface{}) (err error)

type flushRequest struct {
	datas    []interface{}
	doneChan chan struct{}
}
