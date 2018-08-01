terraform {
    required_version = ">= 0.11, < 0.12"
}

provider "citrixitm" {
    client_id     = "${var.itm_client_id}"
    client_secret = "${var.itm_client_secret}"
}

resource "citrixitm_dns_app" "website" {
    name             = "Simple Round Robin"
    description      = "DNS routing for the website"
    app_data         = "${file("${path.module}/website.dns.js")}"
    fallback_cname   = "origin.example.com"
}
