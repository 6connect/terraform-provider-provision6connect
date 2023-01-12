
resource "provision6connect_resource" "pvrecord" {
  parent_id = "428964"
  type = "dnsrecord"
  name = "TerraForm Record"
  attrs = {
    "record_host" = "terraform23.6ckubs.com."
    "record_value" = "111.111.111.112"
    "record_type" = "A"
    "record_ttl" = "900"
  }
}

output "pv_record" {
  value = provision6connect_resource.pvrecord
}
