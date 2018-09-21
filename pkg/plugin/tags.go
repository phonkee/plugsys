package plugin

import "strings"

// AppTag is struct tag
type AppTag string

func (p AppTag) PluginID() string {
	parts := strings.Split(string(p), ",")
	return strings.TrimSpace(parts[0])
}

func (p AppTag) Optional() (result bool) {
	parts := strings.Split(string(p), ",")
	if len(parts) < 2 {
		result = false
	} else {
		result = strings.TrimSpace(parts[1]) == "optional"
	}
	return
}