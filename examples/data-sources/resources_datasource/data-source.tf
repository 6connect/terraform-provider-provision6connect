
data "provision6connect_resources" "kubs" {
  search = {
    "type" = "dnszone"
    "limit" = 2
  }
}

output "kubs_resources" {
  value = data.provision6connect_resources.kubs
}
