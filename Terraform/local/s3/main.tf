resource "aws_s3_bucket" "organize-this-local-bucket" {
    bucket = "${var.project_name}-${var.environment}-qr-codes"
}
