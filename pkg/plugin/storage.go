package plugin

import (
	"sync"

	"github.com/blang/semver"
	"github.com/phonkee/plugsys/api"
	"github.com/phonkee/plugsys/pkg/injector"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	InjectPluginTag = "plugin"
)

// NewStorage returns new plugin storage
func NewStorage(logger *zap.Logger) api.PluginStorage {
	return &storage{
		logger: logger,
		apps:   []api.PluginStorageItem{},
		inj:    injector.New(),
	}
}

// apps implements AppStorage
type storage struct {
	logger *zap.Logger
	apps   []api.PluginStorageItem
	mutex  sync.RWMutex
	inj    injector.Injector
}

// Add adds plugin to storage
func (s *storage) Add(plugin api.Plugin) (err error) {
	s.logger.Debug("adding plugin", zap.String("plugin", plugin.ID()))

	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, p := range s.apps {
		if p.Plugin().ID() == plugin.ID() {
			return errors.Wrap(api.ErrPluginAlreadyRegistered, plugin.ID())
		}
	}

	// add plugin to apps
	s.apps = append(s.apps, newStorageItem(plugin))

	return s.inj.Provide(plugin, plugin.ID(), "plugin")
}

// Each iterates over all apps and calls callback
// If callback returns error, iteration is stopped and the error is propagated to the caller.
// If callback returns StopIteration, iteration is stopped and no error is returned
func (s *storage) Each(callback func(app api.Plugin) error) (err error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, p := range s.apps {
		if err = callback(p.Plugin()); err != nil {
			if err == api.StopIteration {
				err = nil
			}
			break
		}
	}

	return
}

// Apps filters apps by given callback. callback must have single argument which is interface.
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

	if err = s.Each(func(p api.Plugin) (errCallback error) {
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
		err = api.ErrNotApplied
	}
	return
}

func (s *storage) Get(ID string) (api.Plugin, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, p := range s.apps {
		if p.Plugin().ID() == ID {
			return p.Plugin(), true
		}
	}

	return nil, false
}

// Exists returns whether plugin is already added
func (s *storage) Exists(ID string) (exists bool) {
	_, exists = s.Get(ID)
	return
}

// Len returns count of apps
func (s *storage) Len() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return len(s.apps)
}

// Inject injects plugins into given struct
func (s *storage) Inject(target interface{}, skipMissing bool) (err error) {
	return s.inj.Inject(target, skipMissing)
}

// Available returns all plugin ids available
func (s *storage) Available() (result []string) {

	result = make([]string, 0, s.Len())

	s.Each(func(plugin api.Plugin) error {
		result = append(result, plugin.ID())
		return nil
	})
	return
}

// Version returns plugin Version
func (s *storage) Version(ID string) (semver.Version, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, p := range s.apps {
		if p.Plugin().ID() == ID {
			return p.Version(), true
		}
	}

	return semver.Version{}, false
}
