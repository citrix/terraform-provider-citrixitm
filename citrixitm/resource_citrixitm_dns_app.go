package citrixitm

import (
	"fmt"
	"log"
	"strconv"

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
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"fallback_cname": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cname": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCitrixITMDnsAppCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*itm.Client)
	opts := itm.NewDnsAppOpts(
		d.Get("name").(string),
		d.Get("description").(string),
		d.Get("fallback_cname").(string),
		d.Get("app_data").(string),
	)
	log.Printf("[DEBUG] DNS app create options: %#v\n", opts)
	app, err := client.DnsApps.Create(&opts)
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
	app, err := c.DnsApps.Get(id)
	if err != nil {
		return fmt.Errorf("Error retrieving app: %s", err)
	}
	d.Set("name", app.Name)
	d.Set("description", app.Description)
	d.Set("fallback_cname", app.FallbackCname)
	d.Set("app_data", app.AppData)
	d.Set("cname", app.AppCname)
	return nil
}

func update(id int, c *itm.Client, d *schema.ResourceData) error {
	if d.HasChange("name") ||
		d.HasChange("description") ||
		d.HasChange("fallback_cname") ||
		d.HasChange("app_data") {
		opts := itm.NewDnsAppOpts(
			d.Get("name").(string),
			d.Get("description").(string),
			d.Get("fallback_cname").(string),
			d.Get("app_data").(string),
		)
		log.Printf("[DEBUG] DNS app update options: %#v\n", opts)
		_, err := c.DnsApps.Update(id, &opts)
		if err != nil {
			return nil
		}
	}
	return read(id, c, d)
}

func delete(id int, c *itm.Client, d *schema.ResourceData) error {
	err := c.DnsApps.Delete(id)
	if err != nil {
		return fmt.Errorf("There was a problem deleting app (id %d): %v", id, err)
	}
	return nil
}
