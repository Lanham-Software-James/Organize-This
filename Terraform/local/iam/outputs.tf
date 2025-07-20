output "access_key" {
  description = "The access key ID of the system user."
  value       = aws_iam_access_key.organize-this-local-system-key.id
}

output "secret_access_key" {
  description = "The access key secret of the system user."
  value       = aws_iam_access_key.organize-this-local-system-key.secret
  sensitive   = true
}
