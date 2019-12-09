package eventbus

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestNewBus(t *testing.T) {
	bus := NewBus(context.TODO())
	c := make(chan interface{})
	go func() {
		select {
		case event := <-c:
			fmt.Println(event)
		}
	}()
	bus.Subscribe("test", c)

	bus.Send("test", "test")
	time.Sleep(2*time.Second)
}
