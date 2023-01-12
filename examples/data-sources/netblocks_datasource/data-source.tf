
data "provision6connect_netblocks" "someblock" {
  search = {
    "id" = "16303"
    
  }
}

output "kubs_sss" {
  value = data.provision6connect_netblocks.someblock
}
