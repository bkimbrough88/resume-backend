data "aws_route53_zone" "zone" {
  name = var.base_domain_name
}

data "aws_iam_policy_document" "lambda_assumer" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRole"]
    principals {
      identifiers = ["lambda.amazonaws.com"]
      type        = "Service"
    }
  }
}

data "aws_iam_policy_document" "lambda_permissions" {
  statement {
    sid    = "ReadWriteTable"
    effect = "Allow"
    actions = [
      "dynamodb:BatchGetItem",
      "dynamodb:DeleteItem",
      "dynamodb:GetItem",
      "dynamodb:Query",
      "dynamodb:Scan",
      "dynamodb:BatchWriteItem",
      "dynamodb:PutItem",
      "dynamodb:UpdateItem"
    ]
    resources = [aws_dynamodb_table.table.arn]
  }
  statement {
    sid    = "LambdaLogs"
    effect = "Allow"
    actions = [
      "logs:CreateLogStream",
      "logs:PutLogEvents",
      "logs:CreateLogGroup"
    ]
    resources = ["*"]
  }
}

data "aws_cloudfront_origin_request_policy" "cors" {
  name = "Managed-CORS-S3Origin"
}

data "archive_file" "zip" {
  output_path = "${var.function_base_path}/resume-backend.zip"
  source_file = "${var.function_base_path}/resume-backend"
  type        = "zip"
}