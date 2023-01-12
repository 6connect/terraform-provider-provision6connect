
resource "provision6connect_dnszone" "tfexample" {
  name = "tfexample.com."
  group_id = "799411"
}

resource "provision6connect_dnsrecord" "pvrecord" {
  zone_id = provision6connect_dnszone.tfexample.id
  name = "TerraForm Record"
  record_host = "terraform.tfexample.com."
  record_value = "111.111.111.112"
  record_type = "A"
  record_ttl = "900"

  depends_on = [
    provision6connect_dnszone.tfexample
  ]
}

data "provision6connect_dnspush" "pushtf" {
  zone_id = provision6connect_dnszone.tfexample.id
  depends_on = [
    provision6connect_dnsrecord.pvrecord
  ]
}

data "provision6connect_dnspushstatus" "pushtfstatus" {
  zone_id = provision6connect_dnszone.tfexample.id
  push_pid = data.provision6connect_dnspush.pushtf.push_pid
  delay = 5000
  depends_on = [
    data.provision6connect_dnspush.pushtf
  ]
}

output "pv_zone" {
  value = provision6connect_dnszone.tfexample
}

output "push_status" {
  value = data.provision6connect_dnspushstatus.pushtfstatus
}


