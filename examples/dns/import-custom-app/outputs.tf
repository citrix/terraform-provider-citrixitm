output "dns_app_id" {
    value = "${citrixitm_dns_app.simple_app.id}"
}

output "dns_app_version" {
    value = "${citrixitm_dns_app.simple_app.version}"
}

output "dns_app_cname" {
    value = "${citrixitm_dns_app.simple_app.cname}"
}
