package dispatcher

import "context"

type BlockPipeline func(ctx context.Context, values map[string]interface{}) error

type BlockPipelines struct {
	Ctx        context.Context
	CancelFunc context.CancelFunc
	Values     map[string]interface{}

	steps   []BlockPipeline
	current int
}

func NewBlockPipelines(parent context.Context) *BlockPipelines {
	cancel, cancelFunc := context.WithCancel(parent)
	return &BlockPipelines{
		Ctx:        cancel,
		CancelFunc: cancelFunc,
		Values:     make(map[string]interface{}),
		steps:      make([]BlockPipeline, 0),
		current:    -1,
	}
}

func (p *BlockPipelines) AddPipeline(pipeline ...BlockPipeline) {
	p.steps = append(p.steps, pipeline...)
}

func (p *BlockPipelines) Run(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return nil
	case <-p.Ctx.Done():
		return nil
	default:
		if p.current < len(p.steps) {
			e := p.steps[p.current+1](p.Ctx, p.Values)
			if e != nil {
				return e
			}
		}
	}
	return nil
}
