package dispatcher

import (
	"context"
	"github.com/sirupsen/logrus"
)

type worker struct {
	pipelines chan Pipelines
	Name      int
}

func (w *worker) Start(ctx context.Context, works chan<- chan Pipelines) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(w.pipelines)
				return
			default:
				works <- w.pipelines
				select {
				case pipelines, ok := <-w.pipelines:
					if !ok {
						return
					}
					execPipelines(ctx, pipelines)
				}
			}
		}
	}()
}

func execPipelines(ctx context.Context, pipelines Pipelines) {
	defer func() {
		if e := recover(); e != nil {
			logrus.Errorf(`worker exec error:%v`, e)
		}
	}()
	err := pipelines.Run(ctx)
	if err != nil {
		logrus.Errorf(`worker exec error:%v`, err)
	}
	logrus.Infof(`任务执行完成`)
}

func newWorker(i int) *worker {
	return &worker{
		pipelines: make(chan Pipelines),
		Name:      i,
	}
}
