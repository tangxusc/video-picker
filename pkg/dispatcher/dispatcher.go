package dispatcher

import (
	"context"
	"github.com/sirupsen/logrus"
)

type Dispatcher struct {
	MaxWorkers int
	workers    chan chan *Job
	jobs       chan *Job
}

type Job struct {
	F          func(ctx context.Context) error
	JobCtx     context.Context
	CancelFunc context.CancelFunc
}

func NewJob(f func(ctx context.Context) error) *Job {
	return &Job{F: f}
}

func NewDispatcher(maxWorkers int) *Dispatcher {
	d := &Dispatcher{
		MaxWorkers: maxWorkers,
		workers:    make(chan chan *Job, maxWorkers),
		jobs:       make(chan *Job),
	}
	if maxWorkers < 1 {
		panic("worker必须至少1个以上")
	}

	return d
}

func (d *Dispatcher) Dispatch(target *Job) {
	//go func() {
	d.jobs <- target
	//}()
}

func (d *Dispatcher) Start(ctx context.Context) {
	for i := 0; i < d.MaxWorkers; i++ {
		newWorker().Start(ctx, d.workers)
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
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
