package eventbus

import (
	"context"
	"sync"
)

type Bus struct {
	topics map[string][]chan interface{}
	ctx context.Context
	mu  sync.Locker
}

func NewBus(ctx context.Context) *Bus {
	return &Bus{
		topics: make(map[string][]chan interface{}),
		ctx: ctx,
		mu:  &sync.Mutex{},
	}
}

func (b *Bus) Subscribe(topic string, ch chan interface{}) {
	select {
	case <-b.ctx.Done():
		panic("bus closed")
	default:
		chans, ok := b.topics[topic]
		b.mu.Lock()
		defer b.mu.Unlock()
		if !ok {
			chans = make([]chan interface{}, 1)
			chans[0] = ch
			b.topics[topic] = chans
		} else {
			chans = append(chans, ch)
		}
	}
}

func (b *Bus) Send(topic string, event interface{}) {
	chans, ok := b.topics[topic]
	if !ok {
		return
	}
	select {
	case <-b.ctx.Done():
		panic("bus closed")
	default:
		for _, c := range chans {
			go func(e interface{}) { c <- e }(event)
		}
	}
}
