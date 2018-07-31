
terraform {
    required_version = ">= 0.11, < 0.12"
}

provider "digitalocean" {
    version = "~> 0.1"
    token = "${var.digitalocean_api_token}"
}

provider "template" {
    version = "~> 0.1"
}

module "ws1" {
    source                  = "../modules/digitalocean/webserver"
    droplet_name            = "webserver1"
    digitalocean_region     = "sfo2"
    digitalocean_ssh_keys   = "${var.digitalocean_ssh_keys}"
}

module "ws2" {
    source                  = "../modules/digitalocean/webserver"
    droplet_name            = "webserver2"
    digitalocean_region     = "nyc3"
    digitalocean_ssh_keys   = "${var.digitalocean_ssh_keys}"
}

resource "digitalocean_firewall" "webservers" {
  name = "weighted-round-robin-webserver-firewall"

  droplet_ids = [
    "${module.ws1.droplet_id}",
    "${module.ws2.droplet_id}"
  ]

  // The inbound port 22 rule isn't really necessary. We won't be allowing
  // SSH connections at all once everything works.
  inbound_rule = [
    {
      protocol           = "tcp"
      port_range         = "22"
      source_addresses   = "${split(",",var.digitalocean_allowed_ips)}"
    },
    {
      protocol           = "tcp"
      port_range         = "80"
      source_addresses   = ["0.0.0.0/0", "::/0"]
    }
  ]
}
