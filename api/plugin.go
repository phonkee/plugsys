package api

import "github.com/blang/semver"

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

type PluginStorage interface {
	Add(Plugin) error
	Each(func(Plugin) error) error
	Filter(callback interface{}) (err error)
	Get(ID string) (Plugin, bool)
	Exists(ID string) bool
	Len() int
	Inject(target interface{}) (err error)
	Version(ID string) (semver.Version, bool)
}

type PluginStorageItem interface {
	Plugin() Plugin
	Version() semver.Version
}
