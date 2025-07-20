output "client_id" {
  value = module.cognito.client_id
}

output "client_secret" {
  value     = module.cognito.client_secret
  sensitive = true
}

output "user_pool_id" {
  value = module.cognito.user_pool_id
}

output "access_key" {
  value = module.iam.system_user_access_key_id
}

output "secretaccess_key" {
  value     = module.iam.system_user_access_key_secret
  sensitive = true
}
