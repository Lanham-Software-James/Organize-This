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

module "cognito" {
  source                = "./cognito"
  project_readable_name = var.project_readable_name
  project_name          = var.project_name
  environment           = var.environment
}
