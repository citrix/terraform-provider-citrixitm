package itm

import (
	"encoding/json"
	"fmt"
	"net/url"
)

const dnsAppsBasePath = "v2/config/applications/dns.json"

type dnsAppOpts struct {
	AppData       string `json:"appData"`
	Description   string `json:"description"`
	FallbackCname string `json:"fallbackCname"`
	Name          string `json:"name"`
	Protocol      string `json:"protocol"`
	Type          string `json:"type"`
}

func NewDnsAppOpts(name string, description string, fallbackCname string, appData string) dnsAppOpts {
	result := dnsAppOpts{
		Name:          name,
		Description:   description,
		FallbackCname: fallbackCname,
		AppData:       appData,
		Type:          "V1_JS",
		Protocol:      "dns",
	}
	return result
}

type DnsApp struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Enabled       bool   `json:"enabled"`
	FallbackCname string `json:"fallbackCname"`
	AppData       string `json:"appData"`
}

type DnsAppsListTestFunc func(*DnsApp) bool

type DnsAppsService interface {
	Create(*dnsAppOpts) (*DnsApp, error)
	Update(int, *dnsAppOpts) (*DnsApp, error)
	Get(int) (*DnsApp, error)
	Delete(int) error
	List(opts ...DnsAppsListTestFunc) ([]DnsApp, error)
}

type DnsAppsServiceImpl struct {
	client *Client
}

func (s *DnsAppsServiceImpl) Create(opts *dnsAppOpts) (*DnsApp, error) {
	jsonOpts, err := json.Marshal(opts)
	if err != nil {
		return nil, err
	}
	qs := &url.Values{
		"publish": []string{
			"false",
		},
	}
	resp, err := s.client.post(dnsAppsBasePath, jsonOpts, qs)
	if 201 != resp.StatusCode {
		return nil, &UnexpectedHTTPStatusError{
			Expected: 201,
			Got:      resp.StatusCode,
		}
	}
	var result DnsApp
	json.Unmarshal(resp.Body, &result)
	return &result, nil
}

func (s *DnsAppsServiceImpl) Update(id int, opts *dnsAppOpts) (*DnsApp, error) {
	jsonOpts, err := json.Marshal(opts)
	if err != nil {
		return nil, err
	}
	qs := &url.Values{
		"publish": []string{
			"false",
		},
	}
	resp, err := s.client.put(getDnsAppPath(id), jsonOpts, qs)
	if 200 != resp.StatusCode {
		return nil, &UnexpectedHTTPStatusError{
			Expected: 200,
			Got:      resp.StatusCode,
		}
	}
	var result DnsApp
	json.Unmarshal(resp.Body, &result)
	return &result, nil
}

func (s *DnsAppsServiceImpl) Get(id int) (*DnsApp, error) {
	var result DnsApp
	resp, err := s.client.get(getDnsAppPath(id))
	if err != nil {
		return nil, err
	}
	if 200 != resp.StatusCode {
		return nil, &UnexpectedHTTPStatusError{
			Expected: 200,
			Got:      resp.StatusCode}
	}
	json.Unmarshal(resp.Body, &result)
	return &result, nil
}

func (s *DnsAppsServiceImpl) Delete(id int) error {
	resp, err := s.client.delete(getDnsAppPath(id))
	if 204 != resp.StatusCode {
		return &UnexpectedHTTPStatusError{
			Expected: 204,
			Got:      resp.StatusCode,
		}
	}
	return err
}

func (s *DnsAppsServiceImpl) List(tests ...DnsAppsListTestFunc) ([]DnsApp, error) {
	resp, err := s.client.get(dnsAppsBasePath)
	if err != nil {
		return nil, err
	}
	var all []DnsApp
	var result []DnsApp
	json.Unmarshal(resp.Body, &all)
	for _, current := range all {
		stillOk := true
		for _, currentTest := range tests {
			stillOk = currentTest(&current)
			if !stillOk {
				break
			}
		}
		if stillOk {
			result = append(result, current)
		}
	}
	return result, nil
}

func getDnsAppPath(id int) string {
	return fmt.Sprintf("%s/%d", dnsAppsBasePath, id)
}
