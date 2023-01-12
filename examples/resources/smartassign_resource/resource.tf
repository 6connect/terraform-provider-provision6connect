
resource "provision6connect_smartassign" "smartas" {
  rir="1918"
  mask=30
  resource_id="799399"
  type="ipv4"
#  tags = ["tag1","tag2"]
}

output "smartas_netblock" {
  value = provision6connect_smartassign.smartas
}
