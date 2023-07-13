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
  rest_api_id       = aws_api_gateway_rest_api.rest_api.id
  stage_name        = "prod"
  stage_description = "Deployed at ${timestamp()}"
}

resource "aws_api_gateway_api_key" "api_key" {
  name    = "api-key"
  enabled = true
}

resource "aws_api_gateway_usage_plan" "usage_plan" {
  name = "usage-plan"

  api_stages {
    api_id = aws_api_gateway_rest_api.rest_api.id
    stage  = aws_api_gateway_deployment.deployment.stage_name
  }
}

resource "aws_api_gateway_usage_plan_key" "usage_plan_key" {
  key_id        = aws_api_gateway_api_key.api_key.id
  key_type      = "API_KEY"
  usage_plan_id = aws_api_gateway_usage_plan.usage_plan.id
}

# Policies

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
