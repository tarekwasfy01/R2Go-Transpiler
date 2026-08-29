package appcore

import (
	"context"
	"fmt"
	"runtime"
	"sync/atomic"
)

// Pool deliberately runs fewer CPU jobs than GOMAXPROCS. With the default
// configuration rtogo leaves two scheduler slots unused by this pool for the
// GUI/event loop and runtime housekeeping. This is not OS-level CPU affinity,
// but it prevents rtogo's own worker pool from saturating every Go P.
type Pool struct {
	jobs   chan func()
	closed atomic.Bool
}

func RecommendedWorkers() int {
	p := runtime.GOMAXPROCS(0)
	if p < 4 {
		p = 4
		runtime.GOMAXPROCS(p)
	}
	return max(1, p-2)
}

func NewPool(workers int) *Pool {
	if workers < 1 {
		workers = 1
	}
	p := &Pool{jobs: make(chan func(), workers*2)}
	for i := 0; i < workers; i++ {
		go func() {
			for job := range p.jobs {
				func() {
					defer func() { _ = recover() }()
					job()
				}()
			}
		}()
	}
	return p
}

func (p *Pool) Submit(job func()) error {
	if job == nil {
		return nil
	}
	if p.closed.Load() {
		return fmt.Errorf("worker pool is closed")
	}
	select {
	case p.jobs <- job:
		return nil
	default:
		return fmt.Errorf("worker queue is busy")
	}
}

func (p *Pool) SubmitContext(ctx context.Context, job func()) error {
	if p.closed.Load() {
		return fmt.Errorf("worker pool is closed")
	}
	select {
	case p.jobs <- job:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (p *Pool) Close() {
	// Never wait for an arbitrary future parser/transpiler implementation here.
	// A buggy engine that ignores context cancellation must not be able to hold
	// the GUI shutdown path hostage. The process will tear down any remaining
	// worker goroutine after the last window closes.
	if p.closed.CompareAndSwap(false, true) {
		close(p.jobs)
	}
}
