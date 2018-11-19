package citrixitm

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDnsApp_importBasic(t *testing.T) {
	resourceName := "citrixitm_dns_app.foo"
	randString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCitrixITMDnsAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCitrixITMDnsAppConfig(randString),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})

}
