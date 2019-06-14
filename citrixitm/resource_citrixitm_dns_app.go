package citrixitm

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/cedexis/go-itm/itm"
	"github.com/hashicorp/terraform/helper/schema"
)

const resourceName = "Citrix ITM DNS app"

func resourceCitrixITMDnsApp() *schema.Resource {
	return &schema.Resource{
		Create: resourceCitrixITMDnsAppCreate,
		Read:   resourceCitrixITMDnsAppRead,
		Update: resourceCitrixITMDnsAppUpdate,
		Delete: resourceCitrixITMDnsAppDelete,

		Schema: map[string]*schema.Schema{
			"app_data": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: resourceCitrixITMDnsAppDiffSuppress,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"fallback_cname": {
				Type:     schema.TypeString,
				Required: true,
			},
			"fallback_ttl": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  20,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceCitrixITMDnsAppCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Creating %s", resourceName)
	client := m.(*itm.Client)
	opts := itm.NewDNSAppOpts(
		d.Get("name").(string),
		d.Get("description").(string),
		d.Get("fallback_cname").(string),
		d.Get("app_data").(string),
	)
	log.Printf("[DEBUG] %s create options:\n%#v", resourceName, opts)
	app, err := client.DNSApps.Create(&opts, true)
	if err != nil {
		return nil
	}
	d.SetId(strconv.Itoa(app.Id))
	log.Printf("[INFO] Created %s with ID %s", resourceName, d.Id())
	return resourceCitrixITMDnsAppRead(d, m)
}

func resourceCitrixITMDnsAppRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Reading %s", resourceName)
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Error converting app id (%s) to an integer: %s", d.Id(), err)
	}
	client := m.(*itm.Client)
	app, err := client.DNSApps.Get(id)
	if err != nil {
		// There was a problem retrieving the app
		log.Printf("[WARN] %s with ID %s not found", resourceName, d.Id())

		// Set the resource ID to "" to indicate that the resource is not present
		d.SetId("")
	} else {
		if app.Enabled {
			d.Set("name", app.Name)
			d.Set("description", app.Description)
			d.Set("fallback_cname", app.FallbackCname)
			d.Set("fallback_ttl", app.FallbackTtl)
			d.Set("app_data", app.AppData)
			d.Set("cname", app.AppCname)
			d.Set("version", app.Version)
			log.Printf("[INFO] Read %s with ID %s", resourceName, d.Id())
		} else {
			// When the app is disabled, Terraform should recreate it with a
			// new ID, which is done by setting the ID to "".
			log.Printf("[WARN] The %s with ID %s is disabled. This means it was likely deleted outside of Terraform. 'terraform apply' will recreate the app if you approve. If you wish Terraform to stop prompting about it, then you may want to remove its configuration.", resourceName, d.Id())
			d.SetId("")
		}
	}
	return nil
}

func resourceCitrixITMDnsAppUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Updating %s", resourceName)
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Error converting app id (%s) to an integer: %s", d.Id(), err)
	}
	client := m.(*itm.Client)
	if d.HasChange("name") ||
		d.HasChange("description") ||
		d.HasChange("fallback_cname") ||
		d.HasChange("fallback_ttl") ||
		d.HasChange("app_data") {
		opts := itm.NewDNSAppOpts(
			d.Get("name").(string),
			d.Get("description").(string),
			d.Get("fallback_cname").(string),
			d.Get("app_data").(string),
		)
		log.Printf("[DEBUG] %s update options:\n%#v", resourceName, opts)
		_, err := client.DNSApps.Update(id, &opts, true)
		if err != nil {
			log.Printf("[WARN] There was an error updating %s with ID %s: %s", resourceName, d.Id(), err)
		}
		log.Printf("[INFO] Updated %s with ID %s", resourceName, d.Id())
	}
	return resourceCitrixITMDnsAppRead(d, m)
}

func resourceCitrixITMDnsAppDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Deleting %s with ID %s", resourceName, d.Id())
	id, _ := strconv.Atoi(d.Id())
	client := m.(*itm.Client)
	err := client.DNSApps.Delete(id)
	if err != nil {
		log.Printf("[WARN] There was an error deleting %s with ID %s: %s", resourceName, d.Id(), err)
	}
	log.Printf("[INFO] Deleted %s with ID %s", resourceName, d.Id())
	return nil
}

func resourceCitrixITMDnsAppDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	return strings.TrimSpace(old) == strings.TrimSpace(new)
}
