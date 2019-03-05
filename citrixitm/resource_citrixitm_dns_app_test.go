package citrixitm

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/cedexis/go-itm/itm"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

var (
	appName        string
	appNameUpdated string
)

func init() {
	resource.AddTestSweepers("citrixitm_dns_app", &resource.Sweeper{
		Name: "citrixitm_dns_app",
		F:    testSweepDnsApps,
	})

}

func testSweepDnsApps(region string) error {
	meta, err := sharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client := meta.(*itm.Client)
	apps, err := client.DNSApps.List()
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Found %d DNS apps to sweep", len(apps))

	for _, app := range apps {
		if strings.HasPrefix(app.Name, "foo-") {
			log.Printf("[INFO] Destroying DNS app %s", app.Name)
			if err := client.DNSApps.Delete(app.Id); err != nil {
				return err
			}
		}
	}

	return nil
}

func TestAccDnsApp_basic(t *testing.T) {
	var app itm.DNSApp
	randString := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	appName = fmt.Sprintf("foo-%s", randString)
	appNameUpdated = fmt.Sprintf("bar-%s", randString)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCitrixITMDnsAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCitrixITMDnsAppConfig(randString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCitrixITMDnsAppExists("citrixitm_dns_app.foo", &app),
					testAccCheckCitrixITMDnsAppAttributes(
						&app,
						&testAccCitrixITMDnsAppExpectedAttributes{
							Name:          appName,
							Description:   "some description",
							AppData:       "// some source",
							FallbackCname: "fallback.foo.com",
						}),
				),
			},
			{
				Config: testAccCheckCitrixITMDnsAppConfig_Rename(randString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCitrixITMDnsAppExists("citrixitm_dns_app.foo", &app),
					testAccCheckCitrixITMDnsAppAttributes(
						&app,
						&testAccCitrixITMDnsAppExpectedAttributes{
							Name:          appNameUpdated,
							Description:   "some description",
							AppData:       "// some source",
							FallbackCname: "fallback.foo.com",
						}),
				),
			},
			{
				Config: testAccCheckCitrixITMDnsAppConfig_ChangeAppData(randString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCitrixITMDnsAppExists("citrixitm_dns_app.foo", &app),
					testAccCheckCitrixITMDnsAppAttributes(
						&app,
						&testAccCitrixITMDnsAppExpectedAttributes{
							Name:          appNameUpdated,
							Description:   "some description",
							AppData:       "// some source foo",
							FallbackCname: "fallback.foo.com",
						}),
				),
			},
			{
				Config: testAccCheckCitrixITMDnsAppConfig_ChangeDescription(randString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCitrixITMDnsAppExists("citrixitm_dns_app.foo", &app),
					testAccCheckCitrixITMDnsAppAttributes(
						&app,
						&testAccCitrixITMDnsAppExpectedAttributes{
							Name:          appNameUpdated,
							Description:   "some description foo",
							AppData:       "// some source foo",
							FallbackCname: "fallback.foo.com",
						}),
				),
			},
			{
				Config: testAccCheckCitrixITMDnsAppConfig_ChangeFallbackCNAME(randString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCitrixITMDnsAppExists("citrixitm_dns_app.foo", &app),
					testAccCheckCitrixITMDnsAppAttributes(
						&app,
						&testAccCitrixITMDnsAppExpectedAttributes{
							Name:          appNameUpdated,
							Description:   "some description foo",
							AppData:       "// some source foo",
							FallbackCname: "fallback.bar.com",
						}),
				),
			},
		},
	})
}

type testAccCitrixITMDnsAppExpectedAttributes struct {
	Name          string
	Description   string
	AppData       string
	FallbackCname string
}

func testAccCheckCitrixITMDnsAppAttributes(got *itm.DNSApp, want *testAccCitrixITMDnsAppExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) (err error) {
		if err = testValues("name", want.Name, got.Name); err != nil {
			return
		}
		if err = testValues("description", want.Description, got.Description); err != nil {
			return
		}
		if err = testValues("fallback CNAME", want.FallbackCname, got.FallbackCname); err != nil {
			return
		}
		if err = testValues("app data", want.AppData, got.AppData); err != nil {
			return
		}
		// Check the app CNAME
		isMatch, _ := regexp.MatchString("\\d-\\d{2}-[0-9a-z]{4}-[0-9a-z]{4}\\.cdx\\.cedexis\\.net", got.AppCname)
		if !isMatch {
			err = fmt.Errorf("The app CNAME is invalid. Got: %s", got.AppCname)
		}
		return
	}
}

func testAccCheckCitrixITMDnsAppExists(key string, app *itm.DNSApp) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		res, ok := s.RootModule().Resources[key]
		if !ok {
			return fmt.Errorf("Not found: %s", key)
		}
		if res.Primary.ID == "" {
			return fmt.Errorf("The app id is not set")
		}
		client := testAccProvider.Meta().(*itm.Client)
		id, err := strconv.Atoi(res.Primary.ID)
		if err != nil {
			return err
		}

		// Query the API to see if an app with the expected id exists.
		gotten, err := client.DNSApps.Get(id)
		if err != nil {
			return err
		}

		// Verify the result as well as possible
		if strconv.Itoa(gotten.Id) != res.Primary.ID {
			return newUnexpectedValueError("App id", res.Primary.ID, strconv.Itoa(gotten.Id))
		}
		*app = *gotten
		return nil
	}
}

func testAccCheckCitrixITMDnsAppConfig(randString string) string {
	return fmt.Sprintf(`
resource "citrixitm_dns_app" "foo" {
  name 				= "foo-%s"
  description		= "some description"
  app_data			= "// some source"
  fallback_cname	= "fallback.foo.com"
}`, randString)
}

func testAccCheckCitrixITMDnsAppConfig_Rename(randString string) string {
	return fmt.Sprintf(`
resource "citrixitm_dns_app" "foo" {
  name 				= "bar-%s"
  description		= "some description"
  app_data			= "// some source"
  fallback_cname	= "fallback.foo.com"
}`, randString)
}

func testAccCheckCitrixITMDnsAppConfig_ChangeAppData(randString string) string {
	return fmt.Sprintf(`
resource "citrixitm_dns_app" "foo" {
  name 				= "bar-%s"
  description		= "some description"
  app_data			= "// some source foo"
  fallback_cname	= "fallback.foo.com"
}`, randString)
}

func testAccCheckCitrixITMDnsAppConfig_ChangeDescription(randString string) string {
	return fmt.Sprintf(`
resource "citrixitm_dns_app" "foo" {
  name 				= "bar-%s"
  description		= "some description foo"
  app_data			= "// some source foo"
  fallback_cname	= "fallback.foo.com"
}`, randString)
}

func testAccCheckCitrixITMDnsAppConfig_ChangeFallbackCNAME(randString string) string {
	return fmt.Sprintf(`
resource "citrixitm_dns_app" "foo" {
  name 				= "bar-%s"
  description		= "some description foo"
  app_data			= "// some source foo"
  fallback_cname	= "fallback.bar.com"
}`, randString)
}

// Test that the app is truly gone
func testAccCheckCitrixITMDnsAppDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*itm.Client)

	for _, r := range s.RootModule().Resources {
		if r.Type == "citrixitm_dns_app" {
			id, err := strconv.Atoi(r.Primary.ID)
			if err != nil {
				return err
			}
			// Check for the existence of an app with this id
			app, err := client.DNSApps.Get(id)
			if err != nil {
				return err
			}
			// The API doesn't really cause apps to be deleted, but the `enabled` flag should now be set to `false`
			if app.Enabled {
				return fmt.Errorf("App %d is still enabled", id)
			}
		}
	}

	return nil
}
