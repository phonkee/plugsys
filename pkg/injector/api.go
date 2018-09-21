package injector

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	// DefaultTag is default identifier in struct tags
	DefaultTag = "injector"
)

var (
	// ErrDependencyAlreadyExist is returned when dependency already exists
	ErrDependencyAlreadyExist = errors.New("dependency already exists")

	ErrDependencyNotFound = errors.New("dependency not found")

	ErrTargetCannotBeSet = errors.New("cannot set values to target")
)

// Option is injector option
type Option func(*injector)

// Injector interface is the heart of injector package
type Injector interface {
	// Provide adds dependency by name
	Provide(dependency interface{}, name string, namespace ...string) error

	// Inject injects all dependencies into given structure
	Inject(target interface{}, skipMissing bool) error

	Remove(name string, namespace ...string) (err error)
}

// New returns new Injector instance
func New(options ...Option) (Injector) {

	// prepare injector
	result := &injector{
		tag:    DefaultTag,
		logger: zap.NewNop(),
		deps:   map[string]interface{}{},
	}

	// apply all options
	for _, opt := range options {
		opt(result)
	}

	return result
}
