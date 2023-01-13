# terraform-provider-provision6connect
Terraform Provider for 6connect ProVision


# provision6connect Provider

Interact with 6connect ProVision.

## Example Usage

```terraform
# Configuration-based authentication
provider "provision6connect" {
  username = "testuser"
  password = "test123"
  host     = "http://provision.example.com"
}

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

output "pv_zone" {
  value = provision6connect_dnszone.tfexample
}
```
