variable "digitalocean_api_token" {
    description = "The API token to use for API requests"
}

variable "digitalocean_ssh_keys" {
    description = "A comma-separated list of SSH key IDs or fingerprints"
}

variable "digitalocean_allowed_ips" {
    description = "The IP addresses from which to allow SSH connections"
}
