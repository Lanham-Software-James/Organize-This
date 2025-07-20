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
  value = module.iam.access_key
}

output "secret_access_key" {
  value     = module.iam.secret_access_key
  sensitive = true
}

output "bucket_name" {
  value = module.s3.bucket_name
}
