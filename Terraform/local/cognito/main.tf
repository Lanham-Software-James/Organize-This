data "aws_ses_domain_identity" "organize-this-local" {
    domain = var.domain
}

resource "aws_ses_identity_policy" "allow_cognito_send" {
  identity = data.aws_ses_domain_identity.organize-this-local.domain
  name     = "AllowCognitoSend"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Service = "cognito-idp.amazonaws.com"
        }
        Action = [
          "SES:SendEmail",
          "SES:SendRawEmail"
        ]
        Resource = data.aws_ses_domain_identity.organize-this-local.arn
      }
    ]
  })
}

resource "aws_cognito_user_pool" "organize-this-local" {
  name                = "${var.project_name}-${var.environment}"
  username_attributes = ["email"]
  auto_verified_attributes = ["email"]

  password_policy {
    minimum_length    = 8
    require_lowercase = true
    require_uppercase = true
    require_numbers   = true
    require_symbols   = true
  }

  verification_message_template {
    default_email_option = "CONFIRM_WITH_CODE"
    email_subject        = "${var.project_readable_name} - Verify your email address!"
    email_message        = "Your code is {####}"
  }

  email_configuration {
    email_sending_account = "DEVELOPER"
    from_email_address    = "noreply-${var.environment}@${var.domain}"
    source_arn            = data.aws_ses_domain_identity.organize-this-local.arn
  }

  schema {
    attribute_data_type      = "String"
    developer_only_attribute = false
    mutable                  = true
    name                     = "email"
    required                 = true

    string_attribute_constraints {
      min_length = 1
      max_length = 256
    }
  }

  schema {
    attribute_data_type      = "String"
    developer_only_attribute = false
    mutable                  = true
    name                     = "birthdate"
    required                 = true

    string_attribute_constraints {
      min_length = 1
      max_length = 256
    }
  }

  schema {
    attribute_data_type      = "String"
    developer_only_attribute = false
    mutable                  = true
    name                     = "family_name"
    required                 = true

    string_attribute_constraints {
      min_length = 1
      max_length = 256
    }
  }

  schema {
    attribute_data_type      = "String"
    developer_only_attribute = false
    mutable                  = true
    name                     = "given_name"
    required                 = true

    string_attribute_constraints {
      min_length = 1
      max_length = 256
    }
  }

}

resource "aws_cognito_user_pool_client" "organize-this-local" {
  name = "${var.project_name}-${var.environment}"

  user_pool_id                  = aws_cognito_user_pool.organize-this-local.id
  generate_secret               = true
  refresh_token_validity        = 90
  prevent_user_existence_errors = "ENABLED"
  explicit_auth_flows = [
    "ALLOW_REFRESH_TOKEN_AUTH",
    "ALLOW_USER_PASSWORD_AUTH"
  ]

}

resource "aws_cognito_user_group" "organize-this-local-admin" {
  name         = "${var.project_name}-${var.environment}-admin"
  user_pool_id = aws_cognito_user_pool.organize-this-local.id
  description  = "Managed by Terraform: ${var.project_readable_name} - Admin User Group"
}

resource "aws_cognito_user_group" "organize-this-local-free" {
  name         = "${var.project_name}-${var.environment}-free"
  user_pool_id = aws_cognito_user_pool.organize-this-local.id
  description  = "Managed by Terraform: ${var.project_readable_name} - Free User Group"
}

resource "aws_cognito_user_group" "organize-this-local-trial" {
  name         = "${var.project_name}-${var.environment}-trial"
  user_pool_id = aws_cognito_user_pool.organize-this-local.id
  description  = "Managed by Terraform: ${var.project_readable_name} - Trial User Group"
}

resource "aws_cognito_user_group" "organize-this-local-paid" {
  name         = "${var.project_name}-${var.environment}-paid"
  user_pool_id = aws_cognito_user_pool.organize-this-local.id
  description  = "Managed by Terraform: ${var.project_readable_name} - Paid User Group"
}
