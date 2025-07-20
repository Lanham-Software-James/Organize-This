output "client_id" {
  description = "The client ID of the Cognito user pool client."
  value       = aws_cognito_user_pool_client.organize-this-local.id
}

output "client_secret" {
  description = "The client secret of the Cognito user pool client."
  value       = aws_cognito_user_pool_client.organize-this-local.client_secret
  sensitive   = true
}

output "user_pool_id" {
  description = "The ID of the Cognito user pool."
  value       = aws_cognito_user_pool.organize-this-local.id
}

output "user_pool_arn" {
  description = "The ARN of the Cognito user pool."
  value       = aws_cognito_user_pool.organize-this-local.arn
}
