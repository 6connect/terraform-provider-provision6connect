
resource "provision6connect_dnsrecord" "pvrecord" {
  zone_id = "428964"
  name = "TerraForm Record 2"
  record_host = "terraform2.6ckubs.com."
  record_value = "111.111.111.111"
  record_type = "A"
  record_ttl = "900"
  
}

output "pv_record" {
  value = provision6connect_dnsrecord.pvrecord
}
