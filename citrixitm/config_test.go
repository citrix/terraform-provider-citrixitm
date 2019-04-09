package citrixitm

import (
	"fmt"
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

func TestUserAgentStringOverride(t *testing.T) {
	baseURL, _ := url.Parse("http://example.com/foo")
	c := newConfig("some id", "some secret", baseURL)
	client, _ := c.Client()
	expectedUserAgentString := fmt.Sprintf("%s/%s (%s)", libraryName, libraryVersion, libraryURL)
	if client.UserAgentString != expectedUserAgentString {
		t.Error(unexpectedValueString("User Agent String", expectedUserAgentString, client.UserAgentString))
	}
}
