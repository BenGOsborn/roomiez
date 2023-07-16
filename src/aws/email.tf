# Scheduler

resource "aws_cloudwatch_event_rule" "weekly_email" {
  name                = "weekly-email"
  schedule_expression = "cron(0 9 ? * SUN *)"
}

resource "aws_cloudwatch_event_target" "schedule" {
  rule = aws_cloudwatch_event_rule.weekly_email.name
  arn  = aws_lambda_function.schedule_email_lambda.arn
}

# SQS

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
