package loader

import (
	"context"
	"sync/atomic"

	"golang.org/x/sync/errgroup"
)

type Loader struct {
	g      *errgroup.Group
	loaded int32
	total  int32
}

func New(ctx context.Context) (*Loader, context.Context) {
	g, ctx := errgroup.WithContext(ctx)
	return &Loader{g, 0, 0}, ctx
}

func Add[T any](ctx context.Context, l *Loader, dest *T, load func(ctx context.Context) (T, error)) {
	l.total++
	l.g.Go(func() error {
		var err error
		*dest, err = load(ctx)
		atomic.AddInt32(&l.loaded, 1)
		return err
	})
}

func (l *Loader) NumLoaded() int {
	return int(atomic.LoadInt32(&l.loaded))
}

func (l *Loader) Total() int {
	return int(l.total)
}

func (l *Loader) Load() error {
	return l.g.Wait()
}
