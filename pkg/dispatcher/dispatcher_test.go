package dispatcher

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestNewDispatcher(t *testing.T) {
	dispatcher := NewDispatcher(2)
	dispatcher.Start(context.TODO())
	go dispatcher.Dispatch(func() error {
		fmt.Println("test1")
		return nil
	})
	go dispatcher.Dispatch(func() error {
		fmt.Println("test2")
		return fmt.Errorf(`test %v`,`test`)
	})
	time.Sleep(time.Second*2)
}
