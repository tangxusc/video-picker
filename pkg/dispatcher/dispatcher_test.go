package dispatcher

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestNewDispatcher(t *testing.T) {
	dispatcher := NewDispatcher(context.TODO(), 2)
	dispatcher.Start()
	pipelines := NewPipelines(context.TODO())
	pipelines.AddPipeline(func(ctx context.Context, values map[string]interface{}) error {
		fmt.Println(`test`)
		return nil
	})
	dispatcher.Dispatch(pipelines)
	time.Sleep(time.Second * 2)
}
