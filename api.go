package plugsys

import (
	"github.com/blang/semver"
	"github.com/pkg/errors"
)

var (
	ErrPluginAlreadyRegistered = errors.New("plugin already registered")
	ErrInvalidCallback         = errors.New("invalid callback")
	ErrInterfaceNotImplemented = errors.New("interface not implemented")
	ErrNotApplied              = errors.New("not applied")
	ErrInvalidDependencyName   = errors.New("invalid dependency name")

	StopIteration = errors.New("iteration stopped")
)

// Identity identifies given object
type Identity interface {
	// ID returns identification of struct
	ID() string
}

// Plugin implementation
type Plugin interface {
	Identity
}

type PluginVersion interface {
	Version() semver.Version
}

// PluginStorage implementation
type PluginStorage interface {

	// Add adds new plugin to storage, if it's already there it returns error
	Add(Plugin) error

	// Each iterates over all plugins and calls callback on each of then
	Each(func(Plugin) error) error

	// Exists checks if plugin with given plugin is already registered
	Exists(ID string) bool

	// Filter iterates over all plugins and calls given callback. Callback must be function returning error
	// and have single argument of any interface. Filter method will then try to match each plugin to this
	// interface, and when plugin satisfies that interface, it is called
	Filter(callback interface{}) (err error)

	// Get returns plugin by plugin id
	Get(ID string) (Plugin, bool)

	// Inject injects all dependencies to given target
	Inject(target interface{}, skipMissing bool) (err error)

	// Len returns count of registered plugins
	Len() int

	// Provide adds additional dependencies to injector.
	Provide(target interface{}, name string, namespace string) error

	// RemoveDependency removes dependency
	RemoveDependency(name string, namespace ...string) (err error)

	// Version returns version for given plugin id, if not implemented it returns default dev version
	Version(ID string) (semver.Version, bool)
}

// PluginStorageItem is internal type to store plugins
type PluginStorageItem interface {

	// Plugin instance
	Plugin() Plugin

	// Version of given plugin
	Version() semver.Version
}
