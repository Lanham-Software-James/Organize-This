variable "aws_profile" {
  description = "The AWS profile to use."
  type        = string
  default     = "lanham-software-james"
}

variable "project_readable_name" {
  description = "The human readable name of the project."
  type        = string
  default     = "WillowSuite Vault"
}

variable "project_name" {
  description = "The name of the project as it will be used in AWS."
  type        = string
  default     = "willowsuite-vault"
}

variable "environment" {
  description = "The environment of the project."
  type        = string
  default     = "local"
}

variable "region" {
  description = "The region of the project."
  type        = string
  default     = "us-east-2"
}

variable "account_id" {
  description = "The account ID of the project."
  type        = string
}
