package httphandler

type Serializer = func(code int, obj any)

type (
	handlerOptions struct {
		serializer Serializer
	}

	handlerOption func(*ginCtx, *handlerOptions)
)

func resolveOptions(ctx *ginCtx, options ...handlerOption) *handlerOptions {
	defaultOpt := &handlerOptions{
		serializer: ctx.JSON,
	}

	for _, opt := range options {
		opt(ctx, defaultOpt)
	}

	return defaultOpt
}

func WithPureJSONSerializer() handlerOption {
	return func(ctx *ginCtx, ho *handlerOptions) {
		ho.serializer = ctx.PureJSON
	}
}
