package types

import (
	"sort"
	"strings"
)

// Meta holds arbitrary key-value metadata
type Meta map[string]string

func ParseMetaFromString(str string) Meta {
	str = strings.TrimSpace(str)
	if len(str) == 0 {
		return nil
	}
	meta := make(Meta)
	kvpairs := strings.Split(str, ",")
	for _, kvp := range kvpairs {
		kv := strings.Split(kvp, "=")
		meta[kv[0]] = kv[1]
	}
	return meta
}

func (m Meta) String() string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var kvstr string
	for _, k := range keys {
		kvstr += k + "=" + m[k] + ","
	}
	if len(kvstr) > 0 {
		return kvstr[:len(kvstr)-1]
	}
	return ""
}

// Equal returns true if the given Meta has the same key-values
// as m.
func (m Meta) Equal(in Meta) bool {
	if len(m) == len(in) {
		for k, v := range in {
			if m[k] == v {
				continue
			}
			return false
		}
		return true
	}
	return false
}
