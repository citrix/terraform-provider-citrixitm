# For demo purposes, the way we distiguish between the webservers is that the
# webapp reports the "name" of the server that it's running on. It knows this
# based on the setting of the environment variable EXAMPLE_SERVER_NAME, which
# is specified in the systemd unit file defining the service. Here the unit
# file is created using a template so that Terraform can inject the name of the
# server when it is created.
data "template_file" "systemd_unit_file" {
    template = "${file("${path.module}/templates/demo-webapp.service")}"
    vars {
        "server_name" = "${var.droplet_name}"
    }
}

resource "digitalocean_droplet" "webserver" {
    name        = "${var.droplet_name}"
    image       = 35446905
    size        = "512mb"
    region      = "${var.digitalocean_region}"
    ssh_keys    = "${split(",",var.digitalocean_ssh_keys)}"

    # Install a systemd unit to run the webapp via Docker.
    provisioner "file" {
        # We use a template in order to inject the server name into the systemd
        # unit configuration
        content     = "${data.template_file.systemd_unit_file.rendered}"
        destination = "/lib/systemd/system/demo-webapp.service"
    }

    # Make a symlink to the service definition file in /etc/systemd/system and
    # start the service.
    provisioner "remote-exec" {
        inline = [
            "ln -s /lib/systemd/system/demo-webapp.service /etc/systemd/system/",
            "systemctl start demo-webapp.service"
        ]
    }
}
