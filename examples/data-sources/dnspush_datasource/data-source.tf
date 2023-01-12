
data "provision6connect_dnspush" "somezone" {
  zone_id = "799412"
}

output "push_pid" {
  value = data.provision6connect_dnspush.somezone
}
