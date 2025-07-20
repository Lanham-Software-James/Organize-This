output "bucket_arn" {
  description = "The name of the S3 bucket."
  value       = aws_s3_bucket.organize-this-local-bucket.arn
}

output "bucket_name" {
  description = "The name of the S3 bucket."
  value       = aws_s3_bucket.organize-this-local-bucket.id
}
