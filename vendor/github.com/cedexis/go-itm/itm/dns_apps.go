package itm

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
)

const dnsAppsBasePath = "v2/config/applications/dns.json"

// DNSAppOpts specifies settings used to create a new Citrix ITM DNS app
type DNSAppOpts struct {
	AppData       string `json:"appData"`
	Description   string `json:"description"`
	FallbackCname string `json:"fallbackCname"`
	Name          string `json:"name"`
	Protocol      string `json:"protocol"`
	Type          string `json:"type"`
}

// NewDNSAppOpts creates a DNSAppOpts struct
func NewDNSAppOpts(name string, description string, fallbackCname string, appData string) DNSAppOpts {
	result := DNSAppOpts{
		Name:          name,
		Description:   description,
		FallbackCname: fallbackCname,
		AppData:       appData,
		Type:          "V1_JS",
		Protocol:      "dns",
	}
	return result
}

// DNSApp species settings of an existing Citrix DNS app
type DNSApp struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Enabled       bool   `json:"enabled"`
	FallbackCname string `json:"fallbackCname"`
	FallbackTtl   int    `json:"ttl"`
	AppData       string `json:"appData"`
	AppCname      string `json:"cname"`
	Version       int    `json:"version"`
}

type dnsAppsListTestFunc func(*DNSApp) bool

type dnsAppsService interface {
	Create(*DNSAppOpts, bool) (*DNSApp, error)
	Update(int, *DNSAppOpts, bool) (*DNSApp, error)
	Get(int) (*DNSApp, error)
	Delete(int) error
	List(opts ...dnsAppsListTestFunc) ([]DNSApp, error)
}

type dnsAppsServiceImpl struct {
	client *Client
}

// Create a DNS app
func (s *dnsAppsServiceImpl) Create(opts *DNSAppOpts, publish bool) (*DNSApp, error) {
	jsonOpts, err := json.Marshal(opts)
	if err != nil {
		return nil, err
	}
	publishVal := "false"
	if publish {
		publishVal = "true"
	}
	qs := &url.Values{
		"publish": []string{
			publishVal,
		},
	}
	resp, err := s.client.post(dnsAppsBasePath, jsonOpts, qs)
	if err != nil {
		log.Printf("Error issuing post request from DNSAppsServiceImpl.Create: %v", err)
		return nil, err
	}
	if 201 != resp.StatusCode {
		return nil, &UnexpectedHTTPStatusError{
			Expected: 201,
			Got:      resp.StatusCode,
		}
	}
	var result DNSApp
	json.Unmarshal(resp.Body, &result)
	return &result, nil
}

// Update a DNS app
func (s *dnsAppsServiceImpl) Update(id int, opts *DNSAppOpts, publish bool) (*DNSApp, error) {
	jsonOpts, err := json.Marshal(opts)
	if err != nil {
		return nil, err
	}
	publishVal := "false"
	if publish {
		publishVal = "true"
	}
	qs := &url.Values{
		"publish": []string{
			publishVal,
		},
	}
	resp, err := s.client.put(getDNSAppPath(id), jsonOpts, qs)
	if err != nil {
		log.Printf("Error issuing put request from DNSAppsServiceImpl.Update: %v", err)
		return nil, err
	}
	if 200 != resp.StatusCode {
		return nil, &UnexpectedHTTPStatusError{
			Expected: 200,
			Got:      resp.StatusCode,
		}
	}
	var result DNSApp
	json.Unmarshal(resp.Body, &result)
	return &result, nil
}

func (s *dnsAppsServiceImpl) Get(id int) (*DNSApp, error) {
	var result DNSApp
	resp, err := s.client.get(getDNSAppPath(id))
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

func (s *dnsAppsServiceImpl) Delete(id int) error {
	resp, err := s.client.delete(getDNSAppPath(id))
	if 204 != resp.StatusCode {
		return &UnexpectedHTTPStatusError{
			Expected: 204,
			Got:      resp.StatusCode,
		}
	}
	return err
}

func (s *dnsAppsServiceImpl) List(tests ...dnsAppsListTestFunc) ([]DNSApp, error) {
	resp, err := s.client.get(dnsAppsBasePath)
	if err != nil {
		return nil, err
	}
	var all []DNSApp
	var result []DNSApp
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

func getDNSAppPath(id int) string {
	return fmt.Sprintf("%s/%d", dnsAppsBasePath, id)
}
