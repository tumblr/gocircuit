// Package limiter schedules job execution while maintaining an upper limit on concurrency
package limiter

import (
	"sync"
)

// Limiter schedules go routines for execution, while ensuring that no more than a
// pre-set limit run at any time.
type Limiter struct {
	ch chan struct{}
	wg sync.WaitGroup
}

func New(m int) *Limiter {
	return (&Limiter{}).Init(m)
}

func (l *Limiter) Init(m int) *Limiter {
	l.ch = make(chan struct{}, m)
	for i := 0; i < m; i++ {
		l.ch <- struct{}{}
	}
	return l
}

func (l *Limiter) Open() {
	// Take an execution permit
	<-l.ch
	l.wg.Add(1)
}

func (l *Limiter) Close() {
	// Replace the execution permit
	l.ch <- struct{}{}
	l.wg.Done()
}

func (l *Limiter) Go(f func()) {
	l.Open()
	go func() {
		f()
		l.Close()
	}()
}

func (l *Limiter) Throttle(f func()) {
	for {
		l.Go(f)
	}
}

func (l *Limiter) Wait() {
	l.wg.Wait()
}
