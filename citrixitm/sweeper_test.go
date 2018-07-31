package citrixitm

import (
	"fmt"
	"net/url"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

func sharedConfigForRegion(region string) (interface{}, error) {
	if os.Getenv("ITM_CLIENT_ID") == "" {
		return nil, fmt.Errorf("empty ITM_CLIENT_ID")
	}

	if os.Getenv("ITM_CLIENT_SECRET") == "" {
		return nil, fmt.Errorf("empty ITM_CLIENT_SECRET")
	}

	if os.Getenv("ITM_BASE_URL") == "" {
		return nil, fmt.Errorf("empty ITM_BASE_URL")
	}

	baseURL, _ := url.Parse(os.Getenv("ITM_BASE_URL"))
	config := newConfig(os.Getenv("ITM_CLIENT_ID"), os.Getenv("ITM_CLIENT_SECRET"), baseURL)

	client, err := config.Client()
	if err != nil {
		return nil, fmt.Errorf("error getting Citrix ITM client")
	}

	return client, nil
}
