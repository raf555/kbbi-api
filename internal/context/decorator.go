package context

import "context"

type ContextDecoratorFn func(ctx context.Context) context.Context
