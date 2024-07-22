terraform {
  required_providers {
    objectscale = {
      source = "registry.terraform.io/dell/objectscale"
    }
  }
}

provider "objectscale" {
  username = "root"
  password = "Password123@"
  endpoint = "https://10.225.108.189:443"
  insecure = true
}

data "objectscale_account" "all" {
}

output "powerscale_account_all" {
  value = data.objectscale_account.all
}
