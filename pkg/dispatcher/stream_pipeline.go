package dispatcher

import "context"

type StreamPipeline func(ctx context.Context, in <-chan interface{}, out chan<- interface{}) error
type StreamChan chan interface{}

type StreamPipelines struct {
	Ctx        context.Context
	CancelFunc context.CancelFunc

	steps   []StreamPipeline
	inChans []StreamChan
}

func (p *StreamPipelines) AddPipeline(pipeline StreamPipeline, in StreamChan) {
	p.steps = append(p.steps, pipeline)
	p.inChans = append(p.inChans, in)
}

func (p *StreamPipelines) Run(ctx context.Context) error {
	select {
	case <-ctx.Done():
		p.CancelFunc()
		return nil
	case <-p.Ctx.Done():
		return nil
	default:
		//-->1 -->2 -->3 -->4
		for i, step := range p.steps {
			in := p.inChans[i]
			var out chan interface{}
			if i+1 < len(p.steps) {
				out = p.inChans[i+1]
			}
			go func() {
				err := step(p.Ctx, in, out)
				if err != nil {
					p.CancelFunc()
				}
			}()
		}
	}
	return nil
}
