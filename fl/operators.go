package fl

const (
	// NoOp represents no operator
	NoOp = ""
)

const (
	// OpNotEqual represents the not equal operator
	OpNotEqual = "!"
	// OpLess ...
	OpLess = "lt"
	// OpGreater ...
	OpGreater = "gt"
	// OpLessEqual ...
	OpLessEqual = "le"
	// OpGreaterEqual ...
	OpGreaterEqual = "ge"
)

// parseOp parses the input string checking if it contains any operators
// It returns the Op and remainder value or a NoOp and the input string.
func parseOp(in string) (string, string) {
	if in[0] == '!' {
		return OpNotEqual, in[1:]
	}
	if in[2] == ':' {
		op := in[:2]
		switch op {
		case OpNotEqual, OpLess, OpGreater,
			OpLessEqual, OpGreaterEqual:
			return op, in[3:]
		}

	}
	return NoOp, in
}
