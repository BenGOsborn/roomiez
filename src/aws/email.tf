# Scheduler

resource "aws_cloudwatch_event_rule" "weekly_email" {
  name                = "weekly-email"
  schedule_expression = "cron(0 9 ? * SUN *)"
}

resource "aws_cloudwatch_event_target" "schedule" {
  rule = aws_cloudwatch_event_rule.weekly_email.name
  arn  = aws_lambda_function.schedule_email_lambda.arn
}

# Queue

resource "aws_sqs_queue" "email_queue" {
  name                       = "email-queue"
  visibility_timeout_seconds = 120
}

# Lambda

resource "aws_lambda_function" "schedule_email_lambda" {
  function_name    = "schedule-email"
  role             = aws_iam_role.schedule_email_lambda_role.arn
  handler          = "main"
  runtime          = "go1.x"
  timeout          = 120
  filename         = "schedule_email.zip"
  source_code_hash = filebase64sha256("schedule_email.zip")

  environment {
    variables = {
      ENV     = "production"
      TABLE   = aws_dynamodb_table.subscriptions.id
      SQS_URL = aws_sqs_queue.email_queue.id
    }
  }
}

resource "aws_lambda_function" "send_email_lambda" {
  function_name    = "send-email"
  role             = aws_iam_role.send_email_lambda_role.arn
  handler          = "main"
  runtime          = "go1.x"
  timeout          = 120
  filename         = "send_email.zip"
  source_code_hash = filebase64sha256("send_email.zip")

  environment {
    variables = {
      ENV             = "production"
      SECRETS_ARN     = aws_secretsmanager_secret.secrets.arn
      UNSUBSCRIBE_URL = "${aws_api_gateway_deployment.deployment.invoke_url}${aws_api_gateway_resource.subscribe_resource.path}"
    }
  }
}

# Permissions

resource "aws_lambda_event_source_mapping" "send_email_sqs_mapping" {
  event_source_arn = aws_sqs_queue.email_queue.arn
  function_name    = aws_lambda_function.send_email_lambda.function_name
  batch_size       = 1
}

resource "aws_lambda_permission" "allow_cloudwatch" {
  statement_id  = "AllowExecutionFromCloudWatch"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.schedule_email_lambda.arn
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.weekly_email.arn
}

# Roles

resource "aws_iam_role" "schedule_email_lambda_role" {
  name = "schedule-email-lambda-role"

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

resource "aws_iam_role_policy_attachment" "schedule_email_lambda_basic" {
  role       = aws_iam_role.schedule_email_lambda_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_role_policy_attachment" "schedule_email_lambda_dynamo" {
  role       = aws_iam_role.schedule_email_lambda_role.name
  policy_arn = aws_iam_policy.subscriptions_dynamo_policy.arn
}

resource "aws_iam_role_policy_attachment" "schedule_email_lambda_sqs" {
  role       = aws_iam_role.schedule_email_lambda_role.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonSQSFullAccess"
}

resource "aws_iam_role" "send_email_lambda_role" {
  name = "send-email-lambda-role"

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

resource "aws_iam_role_policy_attachment" "send_email_lambda_basic" {
  role       = aws_iam_role.send_email_lambda_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_role_policy_attachment" "send_email_lambda_sqs" {
  role       = aws_iam_role.send_email_lambda_role.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonSQSFullAccess"
}

resource "aws_iam_role_policy_attachment" "send_email_lambda_secrets_manager_policy" {
  role       = aws_iam_role.send_email_lambda_role.name
  policy_arn = aws_iam_policy.secrets_manager_policy.arn
}
