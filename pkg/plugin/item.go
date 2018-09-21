package plugin

import (
	"github.com/blang/semver"
	"github.com/phonkee/plugsys/api"
)

func newStorageItem(plugin api.Plugin) api.PluginStorageItem {
	result := &storageItem{
		plugin: plugin,
	}

	if pv, ok := plugin.(api.PluginVersion); ok {
		result.version = pv.Version()
	} else {
		result.version, _ = semver.Parse("0.1-dev")
	}

	return result
}

// storageItem implements PluginRegistryItem
type storageItem struct {
	plugin  api.Plugin
	version semver.Version
}

func (r *storageItem) Plugin() api.Plugin {
	return r.plugin
}

func (r *storageItem) Version() (result semver.Version) {
	return r.version
}
