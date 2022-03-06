package loader

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/murkland/moreio"
	"golang.org/x/sync/errgroup"
)

type Loader struct {
	g      *errgroup.Group
	cb     Callback
	loaded int32
	total  int32
}

func New(ctx context.Context, cb Callback) (*Loader, context.Context) {
	g, ctx := errgroup.WithContext(ctx)
	return &Loader{g, cb, 0, 0}, ctx
}

func Add[T any](ctx context.Context, l *Loader, path string, dest *T, load func(ctx context.Context, f moreio.File) (T, error)) {
	l.total++
	l.g.Go(func() error {
		var err error
		f, err := moreio.Open(ctx, path)
		if err != nil {
			return fmt.Errorf("%w while loading %s", err, path)
		}
		*dest, err = load(ctx, f)

		i := atomic.AddInt32(&l.loaded, 1)
		l.cb(path, int(i), int(l.total))
		return err
	})
}

func (l *Loader) Load() error {
	return l.g.Wait()
}

type Callback func(path string, i int, n int)
