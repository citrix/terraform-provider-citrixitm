package citrixitm

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/cedexis/go-itm/itm"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCitrixITMDnsApp() *schema.Resource {
	return &schema.Resource{
		Create: resourceCitrixITMDnsAppCreate,
		Read:   withExistingResource(read),
		Update: withExistingResource(update),
		Delete: withExistingResource(delete),

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

func resourceCitrixITMDnsAppCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*itm.Client)
	opts := itm.NewDNSAppOpts(
		d.Get("name").(string),
		d.Get("description").(string),
		d.Get("fallback_cname").(string),
		d.Get("app_data").(string),
	)
	log.Printf("[DEBUG] DNS app create options: %#v", opts)
	app, err := client.DNSApps.Create(&opts, true)
	if err != nil {
		return nil
	}
	d.SetId(fmt.Sprintf("%d", app.Id))
	return read(app.Id, client, d)
}

type ProcessAppFunc func(id int, c *itm.Client, d *schema.ResourceData) error

func withExistingResource(f ProcessAppFunc) func(*schema.ResourceData, interface{}) error {
	return func(d *schema.ResourceData, m interface{}) error {
		id, err := strconv.Atoi(d.Id())
		if err != nil {
			return fmt.Errorf("Invalid app id: %s", d.Id())
		}
		return f(id, m.(*itm.Client), d)
	}
}

func read(id int, c *itm.Client, d *schema.ResourceData) error {
	app, err := c.DNSApps.Get(id)
	log.Printf("[DEBUG] Inside read; app: %#v", app)
	if err != nil {
		return fmt.Errorf("Error retrieving app: %s", err)
	}
	if app.Enabled {
		d.Set("name", app.Name)
		d.Set("description", app.Description)
		d.Set("fallback_cname", app.FallbackCname)
		d.Set("fallback_ttl", app.FallbackTtl)
		d.Set("app_data", app.AppData)
		d.Set("cname", app.AppCname)
		d.Set("version", app.Version)
	} else {
		log.Printf("The app is disabled. This likely means that it was deleted outside of Terraform.")
		d.SetId("")
	}
	return nil
}

func update(id int, c *itm.Client, d *schema.ResourceData) error {
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
		log.Printf("[DEBUG] DNS app update options: %#v", opts)
		_, err := c.DNSApps.Update(id, &opts, true)
		if err != nil {
			return nil
		}
	}
	return read(id, c, d)
}

func delete(id int, c *itm.Client, d *schema.ResourceData) error {
	err := c.DNSApps.Delete(id)
	if err != nil {
		return fmt.Errorf("There was a problem deleting app (id %d): %v", id, err)
	}
	return nil
}

func resourceCitrixITMDnsAppDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	return strings.TrimSpace(old) == strings.TrimSpace(new)
}
