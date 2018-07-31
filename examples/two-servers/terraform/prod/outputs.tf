output "webserver1_ipv4_address" {
  value = "${module.ws1.ipv4_address}"
}

output "webserver2_ipv4_address" {
  value = "${module.ws2.ipv4_address}"
}
