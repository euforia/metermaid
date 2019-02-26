package types

import "strings"

type Meta map[string]string

func (m Meta) String() string {
	var kvstr string
	for k, v := range m {
		kvstr += k + "=" + v + ","
	}
	return kvstr[:len(kvstr)-1]
}

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

func ParseMetaFromString(str string) Meta {
	meta := make(Meta)
	kvpairs := strings.Split(str, ",")
	for _, kvp := range kvpairs {
		kv := strings.Split(kvp, "=")
		meta[kv[0]] = kv[1]
	}
	return meta
}
