#!/bin/bash
echo "Building the Terraform provider plugin binary..."
make build
mkdir -p /root/.terraform.d/plugins/linux_amd64
cp /go/bin/terraform-provider-citrixitm \
/root/.terraform.d/plugins/linux_amd64/
echo "Copied the plugin binary to the user's Terraform plugins directory."
echo "You can now use the plugin inside the docker container after running 'terraform init' from within your root module directory."
