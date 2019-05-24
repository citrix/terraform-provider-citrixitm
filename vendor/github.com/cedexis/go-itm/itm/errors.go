package itm

// UnexpectedHTTPStatusError is an error type that outputs expected vs actual HTTP status
type UnexpectedHTTPStatusError struct {
	Expected int
	Got      int
}

func (e UnexpectedHTTPStatusError) Error() string {
	return unexpectedValueString("HTTP status", e.Expected, e.Got)
}
