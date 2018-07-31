output "dns_app_id" {
    value = "${citrixitm_dns_app.website.id}"
}

output "dns_app_cname" {
    value = "${citrixitm_dns_app.website.cname}"
}
