package plugin

import "go.uber.org/zap"

type Option func(*storage)

// WithTag sets struct fields tag
func WithTag(tag string) Option {
	return func(s *storage) {
		s.tag = tag
	}
}

// WithLogger adds logger to be used
func WithLogger(logger *zap.Logger) Option {
	return func(s *storage) {
		s.logger = logger
	}
}