
data "provision6connect_dnspushstatus" "somezonestatus" {
  zone_id = "799412"
  push_pid = "10329"
}

output "push_status" {
  value = data.provision6connect_dnspushstatus.somezonestatus
}
