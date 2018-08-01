package citrixitm

import (
	"log"
	"net/url"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"citrixitm_dns_app": resourceCitrixITMDnsApp(),
		},

		Schema: map[string]*schema.Schema{
			"client_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ITM_CLIENT_ID", nil),
				Description: "The OAuth client id for the Citrix ITM public API.",
			},
			"client_secret": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ITM_CLIENT_SECRET", nil),
				Description: "The OAuth client secret for the Citrix ITM public API.",
			},
			// The user can control which API instance to use in multiple ways:
			//  * By setting the `base_url` property in their configuration (preferred)
			//  * By setting the ITM_BASE_URL environment variable
			//
			// If neither of these things is done, the prod base URL is used
			// based on the DefaultFunc setting below
			"base_url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ITM_BASE_URL", "https://portal.cedexis.com/api"),
				Description: "The base URL for Citrix ITM API requests",
			},
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	baseURL, _ := url.Parse(d.Get("base_url").(string))
	log.Printf("[INFO] New client base URL: %s", baseURL.String())
	config := newConfig(
		d.Get("client_id").(string),
		d.Get("client_secret").(string),
		baseURL,
	)
	return config.Client()
}
