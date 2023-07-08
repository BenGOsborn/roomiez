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
}

# API resources and methods

resource "aws_api_gateway_resource" "rentals_resource" {
  rest_api_id = aws_api_gateway_rest_api.rest_api.id
  parent_id   = aws_api_gateway_rest_api.rest_api.root_resource_id
  path_part   = "rentals"
}

resource "aws_api_gateway_method" "rentals_get_method" {
  rest_api_id   = aws_api_gateway_rest_api.rest_api.id
  resource_id   = aws_api_gateway_resource.rentals_resource.id
  http_method   = "GET"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "lambda_integration" {
  rest_api_id             = aws_api_gateway_rest_api.rest_api.id
  resource_id             = aws_api_gateway_resource.rentals_resource.id
  http_method             = aws_api_gateway_method.rentals_get_method.http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.retrieve_rentals_lambda.invoke_arn
}

resource "aws_api_gateway_method" "rentals_post_method" {
  rest_api_id   = aws_api_gateway_rest_api.rest_api.id
  resource_id   = aws_api_gateway_resource.rentals_resource.id
  http_method   = "POST"
  authorization = "NONE"
}

# Lambda

resource "aws_lambda_function" "retrieve_rentals_lambda" {
  function_name    = "retrieve-rentals"
  role             = aws_iam_role.lambda_role.arn
  handler          = "main"
  runtime          = "go1.x"
  filename         = "retrieve_rentals.zip"
  source_code_hash = filebase64sha256("retrieve_rentals.zip")

  environment {
    variables = {
      SECRETS = aws_secretsmanager_secret.secrets.arn
    }
  }
}

resource "aws_lambda_permission" "apigw_retrieve_rentals_lambda_permission" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.retrieve_rentals_lambda.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.rest_api.execution_arn}/*/${aws_api_gateway_method.rentals_get_method.http_method}${aws_api_gateway_resource.rentals_resource.path}"
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

# SQS

resource "aws_sqs_queue" "rentals_queue" {
  name = "rentals-queue"
}
