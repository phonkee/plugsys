package injector

import "strings"

type Tag string

func (t Tag) split() (result []string) {
	val := strings.TrimSpace(string(t))
	result = make([]string, 0)
	for _, v := range strings.Split(val, ",") {
		if v = strings.TrimSpace(v); v != "" {
			result = append(result, v)
		}
	}
	return
}

func (t Tag) Name() string {
	splitted := t.split()
	if len(splitted) == 0 {
		return ""
	}
	return splitted[0]
}

func (t Tag) Optional() bool {
	splitted := t.split()
	for _, v := range splitted[1:] {
		if v == "optional" {
			return true
		}
	}
	return false
}
