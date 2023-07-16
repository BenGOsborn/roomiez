# API

module "subscribe_resource_cors" {
  source          = "squidfunk/api-gateway-enable-cors/aws"
  version         = "0.3.3"
  api_id          = aws_api_gateway_rest_api.rest_api.id
  api_resource_id = aws_api_gateway_resource.subscribe_resource.id
}

resource "aws_api_gateway_resource" "subscribe_resource" {
  rest_api_id = aws_api_gateway_rest_api.rest_api.id
  parent_id   = aws_api_gateway_rest_api.rest_api.root_resource_id
  path_part   = "subscribe"
}

resource "aws_api_gateway_method" "subscribe_put_method" {
  rest_api_id   = aws_api_gateway_rest_api.rest_api.id
  resource_id   = aws_api_gateway_resource.subscribe_resource.id
  http_method   = "PUT"
  authorization = "NONE"
}

resource "aws_api_gateway_method" "subscribe_get_method" {
  rest_api_id   = aws_api_gateway_rest_api.rest_api.id
  resource_id   = aws_api_gateway_resource.subscribe_resource.id
  http_method   = "GET"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "subscribe_lambda_integration" {
  rest_api_id             = aws_api_gateway_rest_api.rest_api.id
  resource_id             = aws_api_gateway_resource.subscribe_resource.id
  http_method             = aws_api_gateway_method.subscribe_put_method.http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.subscribe_lambda.invoke_arn
}

resource "aws_api_gateway_integration" "unsubscribe_lambda_integration" {
  rest_api_id             = aws_api_gateway_rest_api.rest_api.id
  resource_id             = aws_api_gateway_resource.subscribe_resource.id
  http_method             = aws_api_gateway_method.subscribe_get_method.http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.unsubscribe_lambda.invoke_arn
}

# Lambda

resource "aws_lambda_function" "subscribe_lambda" {
  function_name    = "subscribe"
  role             = aws_iam_role.subscribe_lambda_role.arn
  handler          = "main"
  runtime          = "go1.x"
  timeout          = 30
  filename         = "subscribe.zip"
  source_code_hash = filebase64sha256("subscribe.zip")

  environment {
    variables = {
      SECRETS_ARN = aws_secretsmanager_secret.secrets.arn
      ENV         = "production"
      TABLE       = aws_dynamodb_table.subscriptions.id
    }
  }
}

resource "aws_lambda_function" "unsubscribe_lambda" {
  function_name    = "unsubscribe"
  role             = aws_iam_role.unsubscribe_lambda_role.arn
  handler          = "main"
  runtime          = "go1.x"
  timeout          = 30
  filename         = "unsubscribe.zip"
  source_code_hash = filebase64sha256("unsubscribe.zip")

  environment {
    variables = {
      ENV   = "production"
      TABLE = aws_dynamodb_table.subscriptions.id
    }
  }
}

# Permissions

resource "aws_lambda_permission" "apigw_subscribe_lambda_permission" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.subscribe_lambda.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.rest_api.execution_arn}/*/${aws_api_gateway_method.subscribe_put_method.http_method}${aws_api_gateway_resource.subscribe_resource.path}"
}

resource "aws_lambda_permission" "apigw_unsubscribe_lambda_permission" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.unsubscribe_lambda.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.rest_api.execution_arn}/*/${aws_api_gateway_method.subscribe_get_method.http_method}${aws_api_gateway_resource.subscribe_resource.path}"
}

# Roles

resource "aws_iam_role" "subscribe_lambda_role" {
  name = "subscribe-lambda-role"

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

resource "aws_iam_role_policy_attachment" "subscribe_lambda_basic" {
  role       = aws_iam_role.subscribe_lambda_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_role_policy_attachment" "subscribe_lambda_secrets_manager" {
  role       = aws_iam_role.subscribe_lambda_role.name
  policy_arn = aws_iam_policy.secrets_manager_policy.arn
}

resource "aws_iam_role_policy_attachment" "subscribe_lambda_dynamo" {
  role       = aws_iam_role.subscribe_lambda_role.name
  policy_arn = aws_iam_policy.subscriptions_dynamo_policy.arn
}

resource "aws_iam_role" "unsubscribe_lambda_role" {
  name = "unsubscribe-lambda-role"

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

resource "aws_iam_role_policy_attachment" "unsubscribe_lambda_basic" {
  role       = aws_iam_role.unsubscribe_lambda_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_role_policy_attachment" "unsubscribe_lambda_dynamo" {
  role       = aws_iam_role.unsubscribe_lambda_role.name
  policy_arn = aws_iam_policy.subscriptions_dynamo_policy.arn
}
