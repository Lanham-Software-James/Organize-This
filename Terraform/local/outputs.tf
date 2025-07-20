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
