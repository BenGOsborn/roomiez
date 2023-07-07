terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
  }
}

provider "aws" {
  region = "ap-southeast-2"
}

# API gateway

resource "aws_api_gateway_rest_api" "rest_api" {
  name        = "roomiez-api"
}

resource "aws_api_gateway_resource" "rentals_resource" {
  rest_api_id = aws_api_gateway_rest_api.rest_api.id
  parent_id   = aws_api_gateway_rest_api.rest_api.root_resource_id
  path_part   = "rentals"
}

# Lambda

resource "aws_lambda_function" "retrieve_rentals_lambda" {
  function_name = "retrieve-rentals"
  role          = aws_iam_role.lambda_role.arn
  handler       = "retrieve_rentals"
  runtime = "go1.x"
  filename      = "retrieve_rentals.zip"

  environment {
    variables = {
      KEY1 = "VALUE1"
      KEY2 = "VALUE2"
    }
  }
}