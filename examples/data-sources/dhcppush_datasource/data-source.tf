
data "provision6connect_dhcppush" "somegroup" {
  group_id = "799412"
}

output "push_pid" {
  value = data.provision6connect_dhcppush.somegroup
}
