
data "provision6connect_dhcppushstatus" "somegroupstatus" {
  group_id = "1111"
  push_pid = "123123"
}

output "push_status" {
  value = data.provision6connect_dhcppushstatus.somegroupstatus
}
