package hw05parallelexecution

import (
	"sync"
)

type errorCounter struct {
	limit, cnt int
	mu         sync.Mutex
	termCh     chan<- struct{}
	isFail     bool
}

func (ec *errorCounter) inc() {
	ec.mu.Lock()
	defer ec.mu.Unlock()

	ec.cnt++
	if ec.cnt == ec.limit {
		ec.termCh <- struct{}{}
		close(ec.termCh)
		ec.isFail = true
	}
}

func (ec *errorCounter) isFailed() bool {
	ec.mu.Lock()
	defer ec.mu.Unlock()

	return ec.isFail
}

func newErrorCounter(limit int, termCh chan<- struct{}) *errorCounter {
	return &errorCounter{
		limit:  limit,
		mu:     sync.Mutex{},
		termCh: termCh,
	}
}
