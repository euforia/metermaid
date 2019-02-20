package fl

import (
	"fmt"
	"strconv"
	"time"
)

func MatchInt64(value int64, filter Filter) bool {
	switch filter.Operator {
	case NoOp:
		for _, val := range filter.Values {
			if i, err := strconv.ParseInt(val, 10, 64); err == nil {
				if value == i {
					return true
				}
				continue
			}
		}
		return false
	case OpNotEqual:
		for _, val := range filter.Values {
			if i, err := strconv.ParseInt(val, 10, 64); err == nil {
				if value != i {
					return true
				}
				continue
			}
		}
		return false
	case OpLess:
		for _, val := range filter.Values {
			if i, err := strconv.ParseInt(val, 10, 64); err == nil {
				if value < i {
					return true
				}
				continue
			}
		}
		return false
	case OpLessEqual:
		for _, val := range filter.Values {
			if i, err := strconv.ParseInt(val, 10, 64); err == nil {
				if value <= i {
					return true
				}
				continue
			}
		}
		return false
	case OpGreater:
		for _, val := range filter.Values {
			if i, err := strconv.ParseInt(val, 10, 64); err == nil {
				if value > i {
					return true
				}
				continue
			}
		}
		return false
	case OpGreaterEqual:
		for _, val := range filter.Values {
			if i, err := strconv.ParseInt(val, 10, 64); err == nil {
				if value >= i {
					return true
				}
				continue
			}
		}
		return false
	}

	return false
}

func MatchTime(inTime time.Time, filter Filter) bool {
	switch filter.Operator {
	case NoOp:
		rfc3339 := inTime.Format(time.RFC3339)
		for _, val := range filter.Values {
			if rfc3339 == val {
				return true
			}
		}
		return false
	case OpLess:
		for _, val := range filter.Values {
			if t, err := time.Parse(time.RFC3339, val); err == nil {
				if inTime.UnixNano() < t.UnixNano() {
					return true
				}
			}
		}
		return false
	case OpLessEqual:
		for _, val := range filter.Values {
			if t, err := time.Parse(time.RFC3339, val); err == nil {
				fmt.Println(inTime.UnixNano(), t.UnixNano())
				if inTime.UnixNano() <= t.UnixNano() {
					return true
				}
			}
		}
		return false
	case OpGreater:
		for _, val := range filter.Values {
			if t, err := time.Parse(time.RFC3339, val); err == nil {
				if inTime.UnixNano() > t.UnixNano() {
					return true
				}
			}
		}
		return false
	case OpGreaterEqual:
		for _, val := range filter.Values {
			if t, err := time.Parse(time.RFC3339, val); err == nil {
				if inTime.UnixNano() >= t.UnixNano() {
					return true
				}
			}
		}
		return false
	}

	return false
}

func MatchString(val string, filter Filter) bool {
	switch filter.Operator {
	case NoOp:
		// ORs
		for _, fval := range filter.Values {
			if fval == val {
				return true
			}
		}
		return false

	case OpNotEqual:
		// all values must not equal
		for _, fval := range filter.Values {
			if val == fval {
				return false
			}
		}
		return true
	}

	// Default
	return false
}
