package cplugin

import "context"

type Command[T any] interface {
	Initialize(ctx context.Context, args []string) error
	Get(context.Context) T
	Destroy(context.Context) error
}
