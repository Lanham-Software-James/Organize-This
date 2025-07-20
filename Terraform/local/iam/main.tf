resource "aws_iam_user" "organize-this-local-system-user" {
    name = "${var.project_name}-${var.environment}-system"
}

resource "aws_iam_access_key" "organize-this-local-system-key" {
    user = aws_iam_user.organize-this-local-system-user.name
}

data "aws_iam_policy_document" "organize-this-local-system-s3-policy" {
    statement {
        effect = "Allow"
        actions = [
            "s3:ReplicateObject",
            "s3:PutObject",
            "s3:GetObject",
            "s3:ListBucket",
            "s3:DeleteObject",
            "s3:HeadObject"
        ]
        resources = [
            "${var.bucket_arn}",
            "${var.bucket_arn}/*"
            ]
    }
}

data "aws_iam_policy_document" "organize-this-local-system-cognito-policy" {
    statement {
        effect = "Allow"
        actions = [
            "cognito-idp:CreateUserPoolClient",
            "cognito-idp:UpdateUserPoolClient",
            "cognito-idp:DescribeUserPoolClient",
            "cognito-idp:AdminInitiateAuth",
            "cognito-idp:AdminUserGlobalSignOut",
        ]
        resources = [var.user_pool_arn]
    }
}

resource "aws_iam_user_policy" "organize-this-local-system-s3-policy" {
    name = "${var.project_name}-${var.environment}-system-s3-policy"
    user = aws_iam_user.organize-this-local-system-user.name
    policy = data.aws_iam_policy_document.organize-this-local-system-s3-policy.json
}

resource "aws_iam_user_policy" "organize-this-local-system-cognito-policy" {
    name = "${var.project_name}-${var.environment}-system-cognito-policy"
    user = aws_iam_user.organize-this-local-system-user.name
    policy = data.aws_iam_policy_document.organize-this-local-system-cognito-policy.json
}
