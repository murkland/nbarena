package loader

import (
	"context"

	"golang.org/x/sync/errgroup"
)

type Loader struct {
	g *errgroup.Group
}

func New(ctx context.Context) (*Loader, context.Context) {
	g, ctx := errgroup.WithContext(ctx)
	return &Loader{g}, ctx
}

func Add[T any](ctx context.Context, l *Loader, dest *T, load func(ctx context.Context) (T, error)) {
	l.g.Go(func() error {
		var err error
		*dest, err = load(ctx)
		return err
	})
}

func (l *Loader) Load() error {
	return l.g.Wait()
}
