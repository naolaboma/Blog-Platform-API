package worker

import (
	"Blog-API/internal/domain"
	"context"
	"log"
	"sync"
)

type Pool struct {
	jobs    chan domain.Job
	wg      sync.WaitGroup
	workers int
}

func NewPool(workers int, queueSize int) domain.WorkerPool {
	return &Pool{
		jobs:    make(chan domain.Job, queueSize),
		workers: workers,
	}
}
func (p *Pool) Start() {
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go func(workerID int) {
			defer p.wg.Done()
			log.Printf("Worker %d starting", workerID)
			for job := range p.jobs {
				func() {
					defer func() {
						if r := recover(); r != nil {
							log.Printf("Worker %d recovered from panic: %v", workerID, r)
						}
					}()
					if err := job.Run(context.Background()); err != nil {
						log.Printf("Worker %d failed to run job: %v", workerID, err)
					}
				}()
			}
			log.Printf("Worker %d stopping", workerID)
		}(i)
	}
}
func (p *Pool) Submit(job domain.Job) {
	p.jobs <- job
}
func (p *Pool) Shutdown() {
	log.Printf("Worker pool shutting down...")
	close(p.jobs)
	p.wg.Wait()
	log.Printf("Worker pool stopped")
}
