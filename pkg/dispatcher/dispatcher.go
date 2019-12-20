package dispatcher

import (
	"context"
	"github.com/sirupsen/logrus"
)

type Dispatcher struct {
	MaxWorkers    int
	workers       chan chan Pipelines
	pipelinesChan chan Pipelines
	ctx           context.Context
}

func NewDispatcher(ctx context.Context, maxWorkers int) *Dispatcher {
	d := &Dispatcher{
		MaxWorkers:    maxWorkers,
		workers:       make(chan chan Pipelines, maxWorkers),
		pipelinesChan: make(chan Pipelines),
		ctx:           ctx,
	}
	if maxWorkers < 1 {
		panic("worker必须至少1个以上")
	}

	return d
}

func (d *Dispatcher) Dispatch(target Pipelines) {
	select {
	case <-d.ctx.Done():
		close(d.pipelinesChan)
	case d.pipelinesChan <- target:

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
				close(d.pipelinesChan)
				return
			case target, ok := <-d.pipelinesChan:
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
