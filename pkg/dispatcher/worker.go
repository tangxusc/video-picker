package dispatcher

import (
	"context"
	"github.com/sirupsen/logrus"
)

type worker struct {
	jobBuf chan *Job
	Name   int
}

func (w *worker) Start(ctx context.Context, works chan<- chan *Job) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(w.jobBuf)
				return
			default:
				works <- w.jobBuf
				select {
				case job, ok := <-w.jobBuf:
					if !ok {
						return
					}
					execJob(ctx, job)
				}
			}
		}
	}()
}

func execJob(ctx context.Context, job *Job) {
	defer func() {
		if e := recover(); e != nil {
			logrus.Errorf(`worker exec error:%v`, e)
		}
	}()
	subCtx, cancelFunc := context.WithCancel(ctx)
	job.JobCtx = subCtx
	job.CancelFunc = cancelFunc
	err := job.F(subCtx)
	if err != nil {
		logrus.Errorf(`worker exec error:%v`, err)
	}
	logrus.Infof(`任务执行完成`)
}

func newWorker(i int) *worker {
	return &worker{
		jobBuf: make(chan *Job),
		Name:   i,
	}
}
