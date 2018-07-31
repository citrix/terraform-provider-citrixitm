output "droplet_id" {
    value = "${digitalocean_droplet.webserver.id}"
}

output "ipv4_address" {
    value = "${digitalocean_droplet.webserver.ipv4_address}"
}
