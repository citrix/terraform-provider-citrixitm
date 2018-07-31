variable "droplet_name" {
    description = "The name assigned to the new droplet"
}

variable "digitalocean_region" {
    description = "The DigitalOcean region for the resource"
}

# We won't need SSH keys in production. The goal is to not need shell access.
variable "digitalocean_ssh_keys" {
    description = "A comma-separated list of SSH key IDs or fingerprints"
}
