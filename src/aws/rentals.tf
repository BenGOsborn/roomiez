# API

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

resource "aws_api_gateway_integration" "rentals_get_lambda_integration" {
  rest_api_id             = aws_api_gateway_rest_api.rest_api.id
  resource_id             = aws_api_gateway_resource.rentals_resource.id
  http_method             = aws_api_gateway_method.rentals_get_method.http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.retrieve_rentals_lambda.invoke_arn
}

resource "aws_api_gateway_method" "rentals_post_method" {
  rest_api_id      = aws_api_gateway_rest_api.rest_api.id
  resource_id      = aws_api_gateway_resource.rentals_resource.id
  http_method      = "POST"
  authorization    = "NONE"
  api_key_required = true
}

resource "aws_api_gateway_integration" "rentals_post_lambda_integration" {
  rest_api_id             = aws_api_gateway_rest_api.rest_api.id
  resource_id             = aws_api_gateway_resource.rentals_resource.id
  http_method             = aws_api_gateway_method.rentals_post_method.http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.process_rental_lambda.invoke_arn
}

# Lambda

resource "aws_lambda_function" "retrieve_rentals_lambda" {
  function_name    = "retrieve-rentals"
  role             = aws_iam_role.retrieve_rentals_lambda_role.arn
  handler          = "main"
  runtime          = "go1.x"
  timeout          = "30"
  filename         = "retrieve_rentals.zip"
  source_code_hash = filebase64sha256("retrieve_rentals.zip")

  environment {
    variables = {
      SECRETS_ARN = aws_secretsmanager_secret.secrets.arn
      ENV         = "production"
    }
  }
}

resource "aws_lambda_function" "process_rental_lambda" {
  function_name    = "process-rentals"
  role             = aws_iam_role.process_rental_lambda_role.arn
  handler          = "main"
  runtime          = "go1.x"
  timeout          = "30"
  filename         = "process_rental.zip"
  source_code_hash = filebase64sha256("process_rental.zip")

  environment {
    variables = {
      SECRETS_ARN = aws_secretsmanager_secret.secrets.arn
      ENV         = "production"
    }
  }
}

# Permissions

resource "aws_lambda_permission" "apigw_retrieve_rentals_lambda_permission" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.retrieve_rentals_lambda.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.rest_api.execution_arn}/*/${aws_api_gateway_method.rentals_get_method.http_method}${aws_api_gateway_resource.rentals_resource.path}"
}

resource "aws_lambda_permission" "apigw_process_rental_lambda_permission" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.process_rental_lambda.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.rest_api.execution_arn}/*/${aws_api_gateway_method.rentals_post_method.http_method}${aws_api_gateway_resource.rentals_resource.path}"
}

# Roles

resource "aws_iam_role" "retrieve_rentals_lambda_role" {
  name = "retrieve-rentals-lambda-role"

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

resource "aws_iam_role_policy_attachment" "retrieve_rentals_lambda_basic" {
  role       = aws_iam_role.retrieve_rentals_lambda_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_role_policy_attachment" "retrieve_rentals_lambda_secrets_manager_policy" {
  role       = aws_iam_role.retrieve_rentals_lambda_role.name
  policy_arn = aws_iam_policy.secrets_manager_policy.arn
}

resource "aws_iam_role_policy_attachment" "retrieve_rentals_lambda_location_policy" {
  role       = aws_iam_role.retrieve_rentals_lambda_role.name
  policy_arn = aws_iam_policy.location_policy.arn
}

resource "aws_iam_role" "process_rental_lambda_role" {
  name = "process-rental-lambda-role"

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

resource "aws_iam_role_policy_attachment" "process_rental_lambda_basic" {
  role       = aws_iam_role.process_rental_lambda_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_role_policy_attachment" "process_rental_lambda_secrets_manager" {
  role       = aws_iam_role.process_rental_lambda_role.name
  policy_arn = aws_iam_policy.secrets_manager_policy.arn
}

resource "aws_iam_role_policy_attachment" "process_rental_lambda_location" {
  role       = aws_iam_role.process_rental_lambda_role.name
  policy_arn = aws_iam_policy.location_policy.arn
}
