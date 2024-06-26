data "objectscale_account" "all" {
}

output "powerscale_account_all" {
  value = data.objectscale_account.all
}
