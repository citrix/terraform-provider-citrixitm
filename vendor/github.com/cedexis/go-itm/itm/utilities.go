package itm

import (
	"fmt"
	"net/url"
)

func stringToURL(asString string) *url.URL {
	result, _ := url.Parse(asString)
	return result
}

func unexpectedValueString(label string, expected interface{}, got interface{}) string {
	return fmt.Sprintf("Unexpected value [%s]\nExpected: %v\nGot: %v", label, expected, got)
}

func newUnexpectedValueError(label string, expected interface{}, got interface{}) error {
	return fmt.Errorf(unexpectedValueString(label, expected, got))
}

func testValues(label string, expected interface{}, got interface{}) (err error) {
	if expected != got {
		err = fmt.Errorf(unexpectedValueString(label, expected, got))
	}
	return
}
