terraform {
    required_version = ">= 0.11, < 0.12"
}

provider "citrixitm" {
    client_id = "${var.itm_client_id}"
    client_secrit = "${var.itm_client_secret}"
}

resource "citrixitm_dns_app" "my_app" {
    source = "foo source"
}
