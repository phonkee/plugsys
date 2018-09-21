package injector

import "go.uber.org/zap"

// WithLogger sets logger to injector
func WithLogger(logger *zap.Logger) Option {
	return func(i *injector) {
		i.logger = logger
	}
}

// WithTag sets tag that identifies injected field in structure
func WithTag(tag string) Option {
	return func(i *injector) {
		i.tag = tag
	}
}
