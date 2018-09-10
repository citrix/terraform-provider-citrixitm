---
layout: "citrixitm"
page_title: "Provider: Citrix ITM"
sidebar_current: "docs-citrixitm-index"
description: |-
  The Citrix ITM provider is used to create and maintain Citrix ITM resources. The provider needs to be configured with the proper credentials before it can be used.
---

# Citrix ITM Provider

This provider allows you to manager your [Citrix ITM](https://www.cedexis.com/) infrastructure using Terraform.

Use the navigation on the left to read about the available resources.

## Example Usage

Suppose we have a DNS routing application specified in a file named example.js.

**example.js**

```javascript
function init(config) {
    config.requestProvider('edgecast');
}

function onRequest(request, response) {
    response.respond('edgecast', 'www.example.edgecastcdn.net');
    response.setTTL(600);
}
```

The following Terraform configuration can be used to create and maintain a Citrix ITM application using the contents of that file.

**main.tf**

```hcl
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
```
