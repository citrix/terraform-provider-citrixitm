package itm

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	libraryName            = "go-itm"
	libraryVersion         = "0.0.1"
	libraryURL             = "https://github.com/cedexis/" + libraryName
	defaultBaseURL         = "https://portal.cedexis.com/api/"
	defaultUserAgentString = libraryName + "/" + libraryVersion + " (" + libraryURL + ")"
)

type Client struct {
	httpClient      *http.Client
	BaseURL         *url.URL
	UserAgentString string

	// Services
	DnsApps DnsAppsService
}

type clientOpt func(*Client) error

func BaseURL(baseURL *url.URL) clientOpt {
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

func HTTPClient(httpClient *http.Client) clientOpt {
	return func(c *Client) error {
		c.httpClient = httpClient
		return nil
	}
}

func (c *Client) parseOptions(opts ...clientOpt) error {
	for _, option := range opts {
		err := option(c)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewClient(opts ...clientOpt) (*Client, error) {
	baseURL, _ := url.Parse(defaultBaseURL)
	result := &Client{
		BaseURL: baseURL,
	}
	result.DnsApps = &DnsAppsServiceImpl{client: result}
	if err := result.parseOptions(opts...); err != nil {
		return nil, err
	}
	if nil == result.httpClient {
		result.httpClient = http.DefaultClient
	}
	return result, nil
}

type Response struct {
	StatusCode int
	Body       []byte
}

func (c *Client) get(path string) (*Response, error) {
	relURL, _ := url.Parse(path)
	apiURL := c.BaseURL.ResolveReference(relURL)
	req, err := http.NewRequest("GET", apiURL.String(), nil)
	if err != nil {
		return nil, newRequestError(err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgentString)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, newRequestError(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return &Response{
		StatusCode: resp.StatusCode,
		Body:       body,
	}, nil
}

func (c *Client) post(path string, data []byte, qsParams *url.Values) (*Response, error) {
	relURL, _ := url.Parse(path)
	apiURL := c.BaseURL.ResolveReference(relURL)
	if qsParams != nil {
		apiURL.RawQuery = qsParams.Encode()
	}
	req, err := http.NewRequest("POST", apiURL.String(), bytes.NewBuffer(data))
	if err != nil {
		return nil, newRequestError(err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", c.UserAgentString)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, newRequestError(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return &Response{
		StatusCode: resp.StatusCode,
		Body:       body,
	}, nil
}

func (c *Client) put(path string, data []byte, qsParams *url.Values) (*Response, error) {
	relURL, _ := url.Parse(path)
	apiURL := c.BaseURL.ResolveReference(relURL)
	if qsParams != nil {
		apiURL.RawQuery = qsParams.Encode()
	}
	req, err := http.NewRequest("PUT", apiURL.String(), bytes.NewBuffer(data))
	if err != nil {
		return nil, newRequestError(err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", c.UserAgentString)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, newRequestError(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return &Response{
		StatusCode: resp.StatusCode,
		Body:       body,
	}, nil
}

func (c *Client) delete(path string) (*Response, error) {
	relURL, _ := url.Parse(path)
	apiURL := c.BaseURL.ResolveReference(relURL)
	req, err := http.NewRequest("DELETE", apiURL.String(), nil)
	if err != nil {
		return nil, newRequestError(err)
	}
	req.Header.Set("User-Agent", c.UserAgentString)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, newRequestError(err)
	}
	defer resp.Body.Close()
	return &Response{
		StatusCode: resp.StatusCode,
		Body:       nil,
	}, nil
}
