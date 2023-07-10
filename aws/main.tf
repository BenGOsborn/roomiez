terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
  }
}

module "api-gateway-enable-cors" {
  source          = "squidfunk/api-gateway-enable-cors/aws"
  version         = "0.3.3"
  api_id          = aws_api_gateway_rest_api.rest_api.id
  api_resource_id = aws_api_gateway_resource.rentals_resource.id
}

provider "aws" {
  region = "ap-southeast-2"
}

# Secrets manager

resource "aws_secretsmanager_secret" "secrets" {
  name = "secrets"
}

# API gateway

resource "aws_api_gateway_rest_api" "rest_api" {
  name = "roomiez-api"
}

resource "aws_api_gateway_deployment" "deployment" {
  rest_api_id = aws_api_gateway_rest_api.rest_api.id
  stage_name  = "prod"

  variables {
    deployed_at = timestamp()
  }
}

# Roles

resource "aws_iam_role" "lambda_role" {
  name = "lambda-role"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF
}

data "aws_iam_policy_document" "secrets_manager_policy" {
  statement {
    effect    = "Allow"
    actions   = ["secretsmanager:GetSecretValue"]
    resources = [aws_secretsmanager_secret.secrets.arn]
  }
}

resource "aws_iam_policy" "secrets_manager_policy" {
  name   = "secrets-manager-policy"
  policy = data.aws_iam_policy_document.secrets_manager_policy.json
}

resource "aws_iam_role_policy_attachment" "lambda_secrets_manager_policy" {
  role       = aws_iam_role.lambda_role.name
  policy_arn = aws_iam_policy.secrets_manager_policy.arn
}

data "aws_iam_policy_document" "location_policy" {
  statement {
    effect    = "Allow"
    actions   = ["geo:SearchPlaceIndexForText"]
    resources = ["*"]
  }
}

resource "aws_iam_policy" "location_policy" {
  name   = "location-policy"
  policy = data.aws_iam_policy_document.location_policy.json
}

resource "aws_iam_role_policy_attachment" "lambda_location_policy" {
  role       = aws_iam_role.lambda_role.name
  policy_arn = aws_iam_policy.location_policy.arn
}
