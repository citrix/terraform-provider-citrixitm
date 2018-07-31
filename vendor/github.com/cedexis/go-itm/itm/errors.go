package itm

type UnexpectedHTTPStatusError struct {
	Expected int
	Got      int
}

func (e UnexpectedHTTPStatusError) Error() string {
	return unexpectedValueString("HTTP status", e.Expected, e.Got)
}

type RequestError struct {
	wrappedError error
}

func newRequestError(e error) RequestError {
	return RequestError{
		wrappedError: e,
	}
}

func (e RequestError) Error() string {
	return e.wrappedError.Error()
}
