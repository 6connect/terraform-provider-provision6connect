
resource "provision6connect_dnszone" "tfexample" {
  name = "tfexample.com."
  group_id = "799411"
}

output "pv_zone" {
  value = provision6connect_dnszone.tfexample
}
