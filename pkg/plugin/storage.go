package plugin

import (
	"strings"
	"sync"

	"github.com/blang/semver"
	"github.com/phonkee/plugsys"
	"github.com/phonkee/plugsys/pkg/injector"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	InjectDefaultTag = "plugsys"
	Namespace        = "plugin"
)

// NewStorage returns new plugin storage
func NewStorage(opts ...Option) plugsys.PluginStorage {
	result := &storage{
		logger:  zap.NewNop(),
		plugins: []plugsys.PluginStorageItem{},
		tag:     InjectDefaultTag,
	}

	for _, opt := range opts {
		opt(result)
	}

	// instantiate injector
	result.inj = injector.New(
		injector.WithTag(result.tag),
		injector.WithLogger(result.logger),
	)

	return result
}

// plugins implements AppStorage
type storage struct {
	logger  *zap.Logger
	plugins []plugsys.PluginStorageItem
	mutex   sync.RWMutex
	inj     injector.Injector
	tag     string
}

// Add adds plugin to storage
func (s *storage) Add(plugin plugsys.Plugin) (err error) {
	s.logger.Debug("adding", zap.String(Namespace, plugin.ID()))

	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, p := range s.plugins {
		if p.Plugin().ID() == plugin.ID() {
			return errors.Wrap(plugsys.ErrPluginAlreadyRegistered, plugin.ID())
		}
	}

	// add plugin to plugins
	s.plugins = append(s.plugins, newStorageItem(plugin))

	return s.inj.Provide(plugin, plugin.ID(), Namespace)
}

// Available returns all plugin ids available
func (s *storage) Available() (result []string) {

	result = make([]string, 0, s.Len())

	_ = s.Each(func(plugin plugsys.Plugin) error {
		result = append(result, plugin.ID())
		return nil
	})
	return
}

// Each iterates over all plugins and calls callback
// If callback returns error, iteration is stopped and the error is propagated to the caller.
// If callback returns StopIteration, iteration is stopped and no error is returned
func (s *storage) Each(callback func(app plugsys.Plugin) error) (err error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, p := range s.plugins {
		if err = callback(p.Plugin()); err != nil {
			if err == plugsys.StopIteration {
				err = nil
			}
			break
		}
	}

	return
}

// Exists returns whether plugin is already added
func (s *storage) Exists(ID string) (exists bool) {
	_, exists = s.Get(ID)
	return
}

// Apps filters plugins by given callback. callback must have single argument which is interface.
// if plugin implements interface, callback is called
func (s *storage) Filter(callback interface{}) (err error) {
	var (
		insp *inspector
	)

	if insp, err = newInspector(callback); err != nil {
		return
	}

	var (
		found bool
	)

	if err = s.Each(func(p plugsys.Plugin) (errCallback error) {
		if !insp.isImplemented(p) {
			return
		}
		// we have found one implementation
		found = true

		if errCallback = insp.call(p); errCallback != nil {
			return
		}
		return
	}); err != nil {
		return
	}

	if !found {
		err = plugsys.ErrNotApplied
	}
	return
}

// Get returns single plugin
func (s *storage) Get(ID string) (plugsys.Plugin, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, p := range s.plugins {
		if p.Plugin().ID() == ID {
			return p.Plugin(), true
		}
	}

	return nil, false
}

// Inject injects plugins into given struct
func (s *storage) Inject(target interface{}, skipMissing bool) (err error) {
	return s.inj.Inject(target, skipMissing)
}

// Len returns count of plugins
func (s *storage) Len() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return len(s.plugins)
}

// Provide additional dependencies
func (s *storage) Provide(dependency interface{}, name string, namespace string) error {
	if name == Namespace {
		return plugsys.ErrInvalidDependencyName
	}

	args := make([]string, 0)
	if ns := strings.TrimSpace(namespace); ns != "" {
		args = append(args, ns)
	}

	return s.inj.Provide(dependency, name, args...)
}

// RemoveDependency removes dependency
func (s *storage) RemoveDependency(name string, namespace ...string) (err error) {
	return s.inj.Remove(name, Namespace)
}

// Version returns plugin Version
func (s *storage) Version(ID string) (semver.Version, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, p := range s.plugins {
		if p.Plugin().ID() == ID {
			return p.Version(), true
		}
	}

	return semver.Version{}, false
}
