package injector

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

		"github.com/pkg/errors"
	"go.uber.org/zap"
)

// injector implements Injector interface
type injector struct {
	logger *zap.Logger
	tag    string
	mutex  sync.RWMutex
	deps   map[string]interface{}
}

func (i *injector) key(name string, namespace ...string) string {
	parts := make([]string, 0, len(namespace)+1)
	parts = append(parts, namespace...)
	parts = append(parts, name)
	return strings.Join(parts, ":")
}

// Provide adds dependency
func (i *injector) Provide(dependency interface{}, name string, namespace ...string) (err error) {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	// add dependency
	i.deps[i.key(name, namespace...)] = dependency

	return
}

func (i *injector) Remove(name string, namespace ...string) (err error) {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	delete(i.deps, i.key(name, namespace...))

	return
}

func (i *injector) Inject(target interface{}, skipMissing bool) (err error) {

	i.mutex.RLock()
	defer i.mutex.RUnlock()

	// get reflect value of target
	value := reflect.ValueOf(target)

	// if we have pointer, get value
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	// check if we can set value
	if !value.CanSet() {
		return ErrTargetCannotBeSet
	}

	var (
		dependency interface{}
		ok         bool
	)

	// iterate over and inject
	for j := 0; j < value.NumField(); j ++ {
		typ := value.Type().Field(j)

		// get tag value
		tagValue := Tag(typ.Tag.Get(i.tag))

		// if we don't have tag that we are interested in, skip
		if tagValue.Name() == "" {
			continue
		}

		// check if we have dependency for given name
		if dependency, ok = i.deps[tagValue.Name()]; !ok {

			// tag is optional
			if tagValue.Optional() || skipMissing {
				continue
			}

			err = errors.Wrapf(ErrDependencyNotFound, "%v not found (%T: %v)", tagValue, target, value.Type().Field(j).Name)
			i.logger.Error(err.Error(), zap.Strings("available", i.Available()))
			return
		}

		// check if we can set given field
		if !value.Field(j).CanSet() {
			i.logger.Debug("found unexported field", zap.String("field", typ.Name))
			return ErrTargetCannotBeSet
		}

		// get reflect value of dependency
		depValue := reflect.ValueOf(dependency)

		// if we have the same type, we can assign value
		if depValue.Type() == typ.Type {
			value.Field(j).Set(depValue)
			continue
		}

		// check if we have interface
		if value.Field(j).Kind() == reflect.Interface {
			if depValue.Type().Implements(typ.Type) {
				value.Field(j).Set(depValue)
			} else {
				return fmt.Errorf("struct %v field %v doesn't implement interface `%v`", depValue.Type(), value.Type().Field(j).Name, typ.Type)
			}
			continue
		}

		// otherwise error
		return fmt.Errorf("invalid inject: field `%v %v` is not %v", typ.Name, typ.Type, depValue.Type())
	}

	return
}

// Available returns available dependencies
func (i *injector ) Available() (result []string) {
	i.mutex.RLock()
	defer i.mutex.RUnlock()
	result = make([]string, 0, len(i.deps))
	for key := range i.deps {
		result = append(result, key)
	}
	return
}