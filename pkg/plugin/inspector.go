package plugin

import (
	"reflect"

	"github.com/phonkee/plugsys"
	"github.com/pkg/errors"
)

// newInspector returns inspector
func newInspector(callback interface{}) (result *inspector, err error) {
	result = &inspector{}
	if err = result.add(callback); err != nil {
		return
	}
	return
}

type inspector struct {
	callback      reflect.Type
	callbackValue reflect.Value
	target        reflect.Type
}

func (i *inspector) add(callback interface{}) (err error) {
	cType := reflect.TypeOf(callback)
	if cType == nil {
		return errors.Wrap(plugsys.ErrInvalidCallback, "cannot accept nil")
	}
	if cType.Kind() != reflect.Func {
		return errors.Wrap(plugsys.ErrInvalidCallback, "must be function")
	}
	if cType.NumIn() != 1 {
		return errors.Wrap(plugsys.ErrInvalidCallback, "single parameter required")
	}

	// assign first parameter
	i.target = cType.In(0)
	i.callback = cType
	i.callbackValue = reflect.ValueOf(callback)
	errorInterface := reflect.TypeOf((*error)(nil)).Elem()

	if cType.NumOut() != 1 || !cType.Out(0).Implements(errorInterface) {
		return errors.Wrap(plugsys.ErrInvalidCallback, "single return value (error) required")
	}

	// we accept only interfaces
	if i.target.Kind() != reflect.Interface {
		return errors.Wrap(plugsys.ErrInvalidCallback, "parameter must be interface")
	}

	return
}

func (i *inspector) isImplemented(target interface{}) bool {
	return reflect.TypeOf(target).Implements(i.target)
}

func (i *inspector) call(target interface{}) (err error) {
	if !i.isImplemented(target) {
		return plugsys.ErrInterfaceNotImplemented
	}

	in := make([]reflect.Value, 1)
	in[0] = reflect.ValueOf(target)
	result := i.callbackValue.Call(in)[0]

	ri := result.Interface()
	if ri == nil {
		return
	}
	return ri.(error)
}
