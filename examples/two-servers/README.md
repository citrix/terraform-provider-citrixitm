# Two Servers

This example demonstrates using Terraform to set up two web servers, which will be used in many of the other examples. We'll use Packer, another tool made by Hashicorp, to generate a virtual machine image with everything needed to spin up fully (or almost fully) provisioned servers. Following the immutable infrastructure paradigm, we never update the configuration of a server. In fact, as you'll see when we set up the firewall, we won't even allow shell access.

Tools and services used in this example:

+ [Packer](https://www.packer.io/) - a tool for automating machine image builds
+ [Terraform](https://www.terraform.io/) - a tool for automating infrastructure deployment
+ [DigitalOcean](https://www.digitalocean.com/) - a cloud computing platform that's developer-friendly and relatively easy to use.
+ [doctl](https://github.com/digitalocean/doctl) - a command line tool for working with DigitalOcean services
