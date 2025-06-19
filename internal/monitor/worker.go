package monitor

import (
	"context"
	"sync"
	"time"
)

// Website represents a website to monitor
type Website struct {
	Name           string
	URL            string
	Method         string
	Headers        map[string]string
	ExpectedStatus int
	CheckContent   string
	Timeout        time.Duration
}

// WorkerPool manages concurrent website checking
type WorkerPool struct {
	workers int
	jobs    chan Website
	results chan CheckResult
	checker *Checker
	wg      sync.WaitGroup
	ctx     context.Context
	cancel  context.CancelFunc
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(workers int, checker *Checker) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	
	return &WorkerPool{
		workers: workers,
		jobs:    make(chan Website, workers*2),
		results: make(chan CheckResult, workers*2),
		checker: checker,
		ctx:     ctx,
		cancel:  cancel,
	}
}

// Start initializes and starts all workers
func (wp *WorkerPool) Start() {
	for i := 0; i < wp.workers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

// worker is the goroutine that processes jobs
func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()
	
	for {
		select {
		case job, ok := <-wp.jobs:
			if !ok {
				return
			}
			
			// Create context with timeout for this specific job
			ctx := wp.ctx
			if job.Timeout > 0 {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(wp.ctx, job.Timeout)
				defer cancel()
			}
			
			result := wp.checker.CheckWebsite(ctx, job)
			
			select {
			case wp.results <- result:
			case <-wp.ctx.Done():
				return
			}
			
		case <-wp.ctx.Done():
			return
		}
	}
}

// Submit adds a website to the job queue
func (wp *WorkerPool) Submit(website Website) {
	select {
	case wp.jobs <- website:
	case <-wp.ctx.Done():
	}
}

// Results returns the results channel
func (wp *WorkerPool) Results() <-chan CheckResult {
	return wp.results
}

// Close shuts down the worker pool
func (wp *WorkerPool) Close() {
	wp.cancel()
	close(wp.jobs)
	wp.wg.Wait()
	close(wp.results)
}