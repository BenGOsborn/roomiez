# API

resource "aws_api_gateway_resource" "rentals_fields_resource" {
  rest_api_id = aws_api_gateway_rest_api.rest_api.id
  parent_id   = aws_api_gateway_resource.rentals_resource.id
  path_part   = "fields"
}

resource "aws_api_gateway_method" "rentals_fields_get_method" {
  rest_api_id   = aws_api_gateway_rest_api.rest_api.id
  resource_id   = aws_api_gateway_resource.rentals_fields_resource.id
  http_method   = "GET"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "rentals_fields_get_lambda_integration" {
  rest_api_id             = aws_api_gateway_rest_api.rest_api.id
  resource_id             = aws_api_gateway_resource.rentals_fields_resource.id
  http_method             = aws_api_gateway_method.rentals_fields_get_method.http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.get_fields_lambda.invoke_arn
}

# Lambda

resource "aws_lambda_function" "get_fields_lambda" {
  function_name    = "get-fields"
  role             = aws_iam_role.lambda_role.arn
  handler          = "main"
  runtime          = "go1.x"
  timeout          = "30"
  filename         = "get_fields.zip"
  source_code_hash = filebase64sha256("get_fields.zip")

  environment {
    variables = {
      SECRETS_ARN = aws_secretsmanager_secret.secrets.arn
      ENV         = "production"
    }
  }
}

# Permissions

resource "aws_lambda_permission" "apigw_get_fields_lambda_permission" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.get_fields_lambda.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.rest_api.execution_arn}/*/${aws_api_gateway_method.rentals_fields_get_method.http_method}${aws_api_gateway_resource.rentals_fields_resource.path}"
}