package dispatcher

import "context"

type Pipelines interface {
	Run(ctx context.Context) error
}
