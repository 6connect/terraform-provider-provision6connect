
data "provision6connect_firstavailableip" "someip" {
  netblock_cidr = "192.168.192.176/28"
}

output "nextip" {
  value = data.provision6connect_firstavailableip.someip
}
