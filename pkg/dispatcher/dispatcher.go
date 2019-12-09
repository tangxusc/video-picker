package dispatcher

import (
	"context"
	"github.com/sirupsen/logrus"
)

type Dispatcher struct {
	MaxWorkers int
	workers    chan chan *Job
	jobs       chan *Job
	ctx        context.Context
}

type Job struct {
	F          func(ctx context.Context) error
	JobCtx     context.Context
	CancelFunc context.CancelFunc
}

func NewJob(f func(ctx context.Context) error) *Job {
	return &Job{F: f}
}

func NewDispatcher(ctx context.Context, maxWorkers int) *Dispatcher {
	d := &Dispatcher{
		MaxWorkers: maxWorkers,
		workers:    make(chan chan *Job, maxWorkers),
		jobs:       make(chan *Job),
		ctx:        ctx,
	}
	if maxWorkers < 1 {
		panic("worker必须至少1个以上")
	}

	return d
}

func (d *Dispatcher) Dispatch(target *Job) {
	select {
	case <-d.ctx.Done():
		close(d.jobs)
	case d.jobs <- target:

	}
}

func (d *Dispatcher) Start() {
	for i := 0; i < d.MaxWorkers; i++ {
		newWorker(i).Start(d.ctx, d.workers)
	}
	go func() {
		for {
			select {
			case <-d.ctx.Done():
				close(d.jobs)
				return
			case target, ok := <-d.jobs:
				if !ok {
					return
				}
				worker := <-d.workers
				logrus.Debugf(`分配任务[%v]到worker[%v]`, target, worker)
				worker <- target
			}
		}
	}()
}
