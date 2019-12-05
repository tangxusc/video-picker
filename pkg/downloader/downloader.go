package downloader

import (
	"context"
	"github.com/sirupsen/logrus"
)

type downloader struct {
	MaxWorkers int
	//工作队列
	workers chan chan string
	jobs    chan string
}

func NewDownloader(maxWorkers int) *downloader {
	d := &downloader{
		MaxWorkers: maxWorkers,
		workers:    make(chan chan string, maxWorkers),
		jobs:       make(chan string, maxWorkers),
	}
	if maxWorkers < 1 {
		panic("worker必须至少1个以上")
	}

	return d
}

func (d *downloader) Dispatch(target string) {
	go func() {
		d.jobs <- target
	}()
}

func (d *downloader) Start(ctx context.Context) {
	for i := 0; i < d.MaxWorkers; i++ {
		newWorker().Start(ctx, d.workers)
	}
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
}
