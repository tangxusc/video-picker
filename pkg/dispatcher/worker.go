package dispatcher

import (
	"context"
	"github.com/sirupsen/logrus"
)

type worker struct {
	jobBuf chan Job
}

func (w *worker) Start(ctx context.Context, works chan<- chan Job) {
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
					execJob(job)
				}
			}
		}
	}()
}

func execJob(job Job) {
	defer func() {
		if e := recover(); e != nil {
			logrus.Errorf(`worker exec error:%v`, e)
		}
	}()
	err := job()
	if err != nil {
		logrus.Errorf(`worker exec error:%v`, err)
	}
}

func newWorker() *worker {
	return &worker{
		jobBuf: make(chan Job),
	}
}