package api

import "errors"

var (
	ErrPluginAlreadyRegistered = errors.New("plugin already registered")

	ErrImproperlyConfigured = errors.New("ImproperlyConfigured")

	ErrInvalidCallback         = errors.New("invalid callback")
	ErrInterfaceNotImplemented = errors.New("interface not implemented")
	ErrNotAllowedDuringNew     = errors.New("not allowed during New")
	ErrNotApplied              = errors.New("not applied")
	ErrNotInterface            = errors.New("not an interface")
)

var (
	StopIteration = errors.New("iteration stopped")
)
