---
layout: "citrixitm"
page_title: "Provider: Citrix ITM"
sidebar_current: "docs-citrixitm-index"
description: |-
  The Citrix ITM provider is used to create and maintain Citrix ITM resources. The provider needs to be configured with the proper credentials before it can be used.
---

# Citrix ITM Provider

This provider allows you to manage [Citrix ITM](https://www.cedexis.com/) infrastructure using Terraform.

Use the navigation on the left to read about the available resources.

## Example Usage

The following set of files can be used to create and maintain a Citrix ITM application via Terraform.

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

**vars.tf**

```hcl
variable "itm_client_id" {
    description = "Client ID for the Citrix ITM API"
}

variable "itm_client_secret" {
    description = "Client secret for the Citrix ITM API"
}
```

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
    name             = "My App"
    description      = "A very simple DNS routing app"
    app_data         = "${file("${path.module}/example.js")}"
    fallback_cname   = "origin.example.com"
}
```

## Authentication

Terraform requires authentication to use the Citrix ITM API. You will need to set up a client ID and secret pair using the Citrix ITM Portal [OAuth Configuration page](https://portal.cedexis.com/ui/api/oauth). Be sure to keep these credentials secure.

## Configuration Reference

The Citrix ITM provider is configured by setting attributes in the provider block.

The following attributes are supported:

* `client_id` - (Required) The client ID to be used for authenticating with the API.

* `client_secret` - (Required) The client secret corresponding to the client ID.

* `base_url` - (Optional) The base URL for Citrix ITM API endpoints. Default: https://portal.cedexis.com/api

  https://portalha.cedexis.com/api can be used to test the provider. **Use with caution!** If a resource is created under one base URL and you later run `terraform apply` with a different base URL, your Terraform state could wind up a mess.
