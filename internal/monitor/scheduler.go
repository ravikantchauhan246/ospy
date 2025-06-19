package monitor

import (
	"context"
	"log"
	"time"
)

// Scheduler manages periodic website checks
type Scheduler struct {
	workerPool *WorkerPool
	websites   []Website
	interval   time.Duration
	ticker     *time.Ticker
	ctx        context.Context
	cancel     context.CancelFunc
}

// NewScheduler creates a new monitoring scheduler
func NewScheduler(workerPool *WorkerPool, websites []Website, interval time.Duration) *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())
	
	return &Scheduler{
		workerPool: workerPool,
		websites:   websites,
		interval:   interval,
		ctx:        ctx,
		cancel:     cancel,
	}
}

// Start begins the scheduled monitoring
func (s *Scheduler) Start() {
	s.ticker = time.NewTicker(s.interval)
	
	// Perform initial check
	go s.runCheck()
	
	// Start periodic checks
	go func() {
		for {
			select {
			case <-s.ticker.C:
				go s.runCheck()
			case <-s.ctx.Done():
				return
			}
		}
	}()
}

// runCheck submits all websites for checking
func (s *Scheduler) runCheck() {
	log.Printf("Starting check of %d websites", len(s.websites))
	
	for _, website := range s.websites {
		s.workerPool.Submit(website)
	}
}

// Stop stops the scheduler
func (s *Scheduler) Stop() {
	if s.ticker != nil {
		s.ticker.Stop()
	}
	s.cancel()
}