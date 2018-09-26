terraform {
    required_version = ">= 0.11, < 0.12"
}

provider "citrixitm" {
    client_id     = "${var.itm_client_id}"
    client_secret = "${var.itm_client_secret}"
}

resource "citrixitm_dns_app" "website" {
    name             = "My Static App"
    description      = "A simple static response app"
    app_data         = "${file("${path.module}/app.js")}"
    fallback_cname   = "origin.example.com"
}
