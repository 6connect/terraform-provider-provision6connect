
resource "provision6connect_netblock" "testblock" {
  cidr="23.23.23.0/28"
  rir="1918"
  allow_sub_assignments=true
  allow_duplicate=true
}

output "test_netblock" {
  value = provision6connect_netblock.testblock
}
