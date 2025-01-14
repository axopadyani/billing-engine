package businesserror

// BusinessError represents a business-specific error with a message and a kind.
type BusinessError struct {
	msg  string
	kind Kind
}

// New creates and returns a new BusinessError.
//
// Parameters:
//   - msg: A string containing the error message.
//   - kind: The Kind of the error, representing its category or type.
//
// Returns:
//
//	A new BusinessError instance.
func New(msg string, kind Kind) *BusinessError {
	return &BusinessError{msg: msg, kind: kind}
}

// Error returns the error message of the BusinessError.
//
// Returns:
//
//	A string containing the error message.
func (e *BusinessError) Error() string {
	return e.msg
}

// Kind returns the kind of the BusinessError.
//
// Returns:
//
//	The Kind of the error, representing its category or type.
func (e *BusinessError) Kind() Kind {
	return e.kind
}
