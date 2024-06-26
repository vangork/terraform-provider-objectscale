terraform {
  required_providers {
    objectscale = {
      source = "registry.terraform.io/dell/objectscale"
    }
  }
}

provider "objectscale" {
  username = "root"
  password = "Password123!"
  endpoint = "https://10.225.108.186:443"
  insecure = true
}
