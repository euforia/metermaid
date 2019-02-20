// Package fl implements a filter language
package fl

import (
	"strings"
)

// ListDelimiter is the delimiter for list values
const ListDelimiter = ","

// ParseValue parses a single value query. For example in key=lt:54, the
// lt:54 should be passed into parse i.e. the value.  It returns
// an Op if it contains one and parses the or's returning the array
func ParseValue(in string) (string, []string) {
	op, field := parseOp(in)
	return op, parseDelimited(field, ListDelimiter)
}

func parseDelimited(in, delimiter string) []string {
	parts := strings.Split(in, delimiter)
	for i := 0; i < len(parts); i++ {
		if parts[i] == "" {
			parts = append(parts[:i], parts[i+1:]...)
			i--
		}
	}
	return parts
}

// Filter holds a parsed Query
type Filter struct {
	Operator string
	Values   []string
}

// Query holds the complet query request
type Query map[string][]Filter

// ParseQuery ...
func ParseQuery(in map[string][]string) Query {
	queries := make(map[string][]Filter)
	for k := range in {
		queries[k] = make([]Filter, 0)
	}

	for k, vals := range in {
		for _, val := range vals {
			var q Filter
			q.Operator, q.Values = ParseValue(val)

			t := queries[k]
			queries[k] = append(t, q)
		}
	}

	return queries
}
