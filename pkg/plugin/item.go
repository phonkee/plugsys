package plugin

import (
	"github.com/blang/semver"
	"github.com/phonkee/plugsys"
)

func newStorageItem(plugin plugsys.Plugin) plugsys.PluginStorageItem {
	result := &storageItem{
		plugin: plugin,
	}

	if pv, ok := plugin.(plugsys.PluginVersion); ok {
		result.version = pv.Version()
	} else {
		result.version, _ = semver.Parse("0.1-dev")
	}

	return result
}

// storageItem implements PluginRegistryItem
type storageItem struct {
	plugin  plugsys.Plugin
	version semver.Version
}

func (r *storageItem) Plugin() plugsys.Plugin {
	return r.plugin
}

func (r *storageItem) Version() (result semver.Version) {
	return r.version
}
