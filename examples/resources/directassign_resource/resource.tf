
resource "provision6connect_directassign" "directas" {
  cidr="192.168.192.176/28"
  resource_id="799399"
}

output "directas_netblock" {
  value = provision6connect_directassign.directas
}
