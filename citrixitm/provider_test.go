package citrixitm

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"citrixitm": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("ITM_CLIENT_ID"); v == "" {
		t.Fatal("ITM_CLIENT_ID must be set for acceptance tests to run properly")
	}
	if v := os.Getenv("ITM_CLIENT_SECRET"); v == "" {
		t.Fatal("ITM_CLIENT_SECRET must be set for acceptance tests to run properly")
	}
}
