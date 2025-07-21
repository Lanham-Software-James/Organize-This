provider "aws" {
  profile = var.aws_profile

  default_tags {
    tags = {
      Project     = var.project_readable_name
      ManagedBy   = "Terraform"
      Environment = var.environment
    }
  }
}

module "s3" {
  source       = "./s3"
  project_name = var.project_name
  environment  = var.environment
}

module "cognito" {
  source                = "./cognito"
  project_readable_name = var.project_readable_name
  project_name          = var.project_name
  environment           = var.environment
  domain                = var.domain
}

module "iam" {
  source        = "./iam"
  project_name  = var.project_name
  environment   = var.environment
  bucket_arn    = module.s3.bucket_arn
  user_pool_arn = module.cognito.user_pool_arn
}
