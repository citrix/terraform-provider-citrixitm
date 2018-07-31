package citrixitm

import (
	"net/url"
	"testing"
)

func TestConfig(t *testing.T) {
	testData := []struct {
		clientId        string
		clientSecret    string
		baseURL         string
		expectedBaseURL string
	}{
		{
			"foo",
			"bar",
			"http://foo.com/api",
			"http://foo.com/api/",
		},
	}
	for _, current := range testData {
		baseURL, _ := url.Parse(current.baseURL)
		c := newConfig(current.clientId, current.clientSecret, baseURL)
		if current.expectedBaseURL != c.BaseURL.String() {
			t.Error(unexpectedValueString("Base URL", current.expectedBaseURL, c.BaseURL.String()))
		}
	}
}
