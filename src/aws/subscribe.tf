# Database

resource "aws_dynamodb_table" "subscriptions" {
  name         = "subscriptions"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "email"

  attribute {
    name = "email"
    type = "S"
  }

  point_in_time_recovery {
    enabled = true
  }
}

# API

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

resource "aws_api_gateway_integration" "subscribe_lambda_integration" {
  rest_api_id             = aws_api_gateway_rest_api.rest_api.id
  resource_id             = aws_api_gateway_resource.subscribe_resource.id
  http_method             = aws_api_gateway_method.subscribe_put_method.http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.subscribe_lambda.invoke_arn
}

# Lambda

resource "aws_lambda_function" "subscribe_lambda" {
  function_name    = "subscribe"
  role             = aws_iam_role.subscribe_lambda_role.arn
  handler          = "main"
  runtime          = "go1.x"
  timeout          = "30"
  filename         = "subscribe.zip"
  source_code_hash = filebase64sha256("subscribe.zip")

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

data "aws_iam_policy_document" "subscriptions_dynamo_policy" {
  statement {
    effect    = "Allow"
    actions   = ["dynamodb:*"]
    resources = [aws_dynamodb_table.subscriptions.arn]
  }
}

resource "aws_iam_policy" "subscriptions_dynamo_policy" {
  name   = "subscriptions-dynamo-policy"
  policy = data.aws_iam_policy_document.subscriptions_dynamo_policy.json
}

