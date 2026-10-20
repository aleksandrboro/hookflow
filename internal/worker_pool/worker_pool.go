package workerpool

import (
	"context"
	"sync"
	"time"

	"github.com/aleksandrboro/hookflow/internal/domain"
	"github.com/aleksandrboro/hookflow/internal/storage"
)

type WorkerPool struct {
	jobChan       chan *domain.Job
	readyJobsChan chan *domain.Job
	workersCount  int
	store         *storage.Storage
	wg            sync.WaitGroup
}

func NewWorkerPool(bufSize, workersCount int, store *storage.Storage) *WorkerPool {
	return &WorkerPool{
		jobChan:       make(chan *domain.Job, bufSize),
		readyJobsChan: make(chan *domain.Job, bufSize),
		workersCount:  workersCount,
		store:         store,
	}
}

func (p *WorkerPool) AddJob(ctx context.Context, job *domain.Job) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case p.jobChan <- job:
		return nil
	}
}

func (p *WorkerPool) StartWorkers(ctx context.Context) {
	go p.startJobsScheduler(ctx)

	for range p.workersCount {
		p.wg.Add(1)
		go p.worker(ctx)
	}
}

func (p *WorkerPool) Wait() {
	p.wg.Wait()
}

func (p *WorkerPool) startJobsScheduler(ctx context.Context) {
	defer close(p.jobChan)

	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-p.jobChan:
			if !ok {
				return
			}

			duration := time.Until(job.RunAt)
			if duration <= 0 {
				p.sendToReady(ctx, job)
				continue
			}

			time.AfterFunc(duration, func() {
				p.sendToReady(ctx, job)
			})
		}
	}
}

func (p *WorkerPool) sendToReady(ctx context.Context, job *domain.Job) {
	select {
	case <-ctx.Done():
		return
	case p.readyJobsChan <- job:
	}
}

func (p *WorkerPool) worker(ctx context.Context) {
	defer p.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-p.readyJobsChan:
			if !ok {
				return
			}
			p.executeJob(ctx, job)
		}
	}
}

func (p *WorkerPool) executeJob(ctx context.Context, job *domain.Job) {
	
}
