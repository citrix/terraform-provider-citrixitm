package itm

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	libraryName            = "go-itm"
	libraryVersion         = "1.0.1"
	libraryURL             = "https://github.com/cedexis/" + libraryName
	defaultBaseURL         = "https://portal.cedexis.com/api/"
	defaultUserAgentString = libraryName + "/" + libraryVersion + " (" + libraryURL + ")"
)

// Client specifies settings for a new ITM client
type Client struct {
	httpClient      *http.Client
	BaseURL         *url.URL
	UserAgentString string

	// Services
	DNSApps dnsAppsService
}

// ClientOpt is a generic type used to specify validated options for creating an ITM client
type ClientOpt func(*Client) error

// BaseURL creates a client option to specify the base URL for use in accessing the API
func BaseURL(baseURL *url.URL) ClientOpt {
	return func(c *Client) error {
		if baseURL != nil {
			if 0 < len(baseURL.Path) && "/" != baseURL.Path[len(baseURL.Path)-1:] {
				baseURL.Path += "/"
			}
			c.BaseURL = baseURL
		}
		return nil
	}
}

// HTTPClient creates a client option used to specify an HTTP client for use by the ITM client being created
func HTTPClient(httpClient *http.Client) ClientOpt {
	return func(c *Client) error {
		c.httpClient = httpClient
		return nil
	}
}

// UserAgentString creates a client option used to specify the user-agent HTTP request header
func UserAgentString(value string) ClientOpt {
	return func(c *Client) error {
		c.UserAgentString = value
		return nil
	}
}

func (c *Client) parseOptions(opts ...ClientOpt) error {
	for _, option := range opts {
		err := option(c)
		if err != nil {
			return err
		}
	}
	return nil
}

// NewClient creates a new ITM client
func NewClient(opts ...ClientOpt) (*Client, error) {
	baseURL, _ := url.Parse(defaultBaseURL)
	result := &Client{
		BaseURL:         baseURL,
		UserAgentString: defaultUserAgentString,
	}
	result.DNSApps = &dnsAppsServiceImpl{client: result}
	if err := result.parseOptions(opts...); err != nil {
		return nil, err
	}
	if nil == result.httpClient {
		result.httpClient = http.DefaultClient
	}
	return result, nil
}

type response struct {
	StatusCode int
	Body       []byte
}

func (c *Client) get(path string) (*response, error) {
	relURL, _ := url.Parse(path)
	apiURL := c.BaseURL.ResolveReference(relURL)
	req, err := http.NewRequest("GET", apiURL.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgentString)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return &response{
		StatusCode: resp.StatusCode,
		Body:       body,
	}, nil
}

func (c *Client) post(path string, data []byte, qsParams *url.Values) (*response, error) {
	relURL, _ := url.Parse(path)
	apiURL := c.BaseURL.ResolveReference(relURL)
	if qsParams != nil {
		apiURL.RawQuery = qsParams.Encode()
	}
	req, err := http.NewRequest("POST", apiURL.String(), bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", c.UserAgentString)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return &response{
		StatusCode: resp.StatusCode,
		Body:       body,
	}, nil
}

func (c *Client) put(path string, data []byte, qsParams *url.Values) (*response, error) {
	relURL, _ := url.Parse(path)
	apiURL := c.BaseURL.ResolveReference(relURL)
	if qsParams != nil {
		apiURL.RawQuery = qsParams.Encode()
	}
	req, err := http.NewRequest("PUT", apiURL.String(), bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", c.UserAgentString)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return &response{
		StatusCode: resp.StatusCode,
		Body:       body,
	}, nil
}

func (c *Client) delete(path string) (*response, error) {
	relURL, _ := url.Parse(path)
	apiURL := c.BaseURL.ResolveReference(relURL)
	req, err := http.NewRequest("DELETE", apiURL.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", c.UserAgentString)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return &response{
		StatusCode: resp.StatusCode,
		Body:       nil,
	}, nil
}
