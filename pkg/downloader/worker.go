package downloader

import "context"

type worker struct {
	jobBuf chan string
}

func (w *worker) Start(ctx context.Context, works chan<- chan string) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				works <- w.jobBuf
				//开始下载
				select {
				case <-ctx.Done():
					close(w.jobBuf)
					return
				case target, ok := <-w.jobBuf:
					if !ok {
						return
					}
					m3u8 := getM3u8(target)
					download(m3u8)
				}
			}
		}
	}()
}

func newWorker() *worker {
	return &worker{
		jobBuf: make(chan string),
	}
}
