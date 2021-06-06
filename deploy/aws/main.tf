resource "aws_iam_role" "lambda_assumer" {
  assume_role_policy = data.aws_iam_policy_document.lambda_assumer.json
  name               = "ResumeBackendLambdaAssumer"
}

resource "aws_iam_policy" "lambda_permissions" {
  name   = "ResumeBackendLambdaPermissions"
  policy = data.aws_iam_policy_document.lambda_permissions.json
}

resource "aws_iam_role_policy_attachment" "lambda_permissions" {
  policy_arn = aws_iam_policy.lambda_permissions.arn
  role       = aws_iam_role.lambda_assumer.name
}

resource "aws_dynamodb_table" "table" {
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "user_id"
  name         = "resume_user"

  attribute {
    name = "user_id"
    type = "S"
  }
}

resource "aws_lambda_function" "resume_backend" {
  filename         = data.archive_file.zip.output_path
  function_name    = "ResumeBackend"
  handler          = "resume-backend"
  role             = aws_iam_role.lambda_assumer.arn
  runtime          = "go1.x"
  source_code_hash = data.archive_file.zip.output_base64sha256
}

resource "aws_lambda_permission" "apigw" {
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.resume_backend.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.api.execution_arn}/*/*/*"
}

resource "aws_apigatewayv2_api" "api" {
  name          = "resume-api"
  protocol_type = "HTTP"
}

resource "aws_apigatewayv2_authorizer" "auth" {
  api_id = aws_apigatewayv2_api.api.id
  authorizer_type = "JWT"
  identity_sources = ["$request.header.Authorization"]
  name = "auth0"

  jwt_configuration {
    audience = [var.auth0_audience]
    issuer = "https://${var.auth0_domain}/"
  }
}

resource "aws_apigatewayv2_integration" "resume_backend" {
  api_id                    = aws_apigatewayv2_api.api.id
  connection_type           = "INTERNET"
  integration_type          = "AWS_PROXY"
  integration_method        = "POST"
  integration_uri           = aws_lambda_function.resume_backend.invoke_arn
}

resource "aws_apigatewayv2_route" "get_user_by_key" {
  api_id             = aws_apigatewayv2_api.api.id
  authorization_type = "NONE"
  operation_name     = "Get User by Key"
  route_key          = "GET /user/{id}"
  target             = "integrations/${aws_apigatewayv2_integration.resume_backend.id}"
}

resource "aws_apigatewayv2_route" "put_user" {
  api_id             = aws_apigatewayv2_api.api.id
  authorizer_id = aws_apigatewayv2_authorizer.auth.id
  authorization_type = "JWT"
  operation_name     = "Put User"
  route_key          = "POST /user"
  target             = "integrations/${aws_apigatewayv2_integration.resume_backend.id}"
}

resource "aws_apigatewayv2_route" "delete_user" {
  api_id             = aws_apigatewayv2_api.api.id
  authorizer_id = aws_apigatewayv2_authorizer.auth.id
  authorization_type = "JWT"
  operation_name     = "Delete User"
  route_key          = "DELETE /user/{id}"
  target             = "integrations/${aws_apigatewayv2_integration.resume_backend.id}"
}

resource "aws_apigatewayv2_deployment" "deployment" {
  api_id      = aws_apigatewayv2_api.api.id
  description = "HTTP API for Resume Backend"
  triggers = {
    redeployment = sha1(join(",", tolist(
      [
        jsonencode(aws_apigatewayv2_integration.resume_backend),
        jsonencode(aws_apigatewayv2_route.get_user_by_key),
        jsonencode(aws_apigatewayv2_route.put_user),
        jsonencode(aws_apigatewayv2_route.delete_user)
      ]
    )))
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_apigatewayv2_stage" "v1" {
  api_id        = aws_apigatewayv2_api.api.id
  deployment_id = aws_apigatewayv2_deployment.deployment.id
  name          = "v1"
}

locals {
  cfn_origin = "resumeBackend"
  domain_name = "resume-api.${var.base_domain_name}"
}

resource "aws_acm_certificate" "cert" {
  provider = aws.usea1

  domain_name       = local.domain_name
  validation_method = "EMAIL"

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_acm_certificate_validation" "validate" {
  provider = aws.usea1

  certificate_arn = aws_acm_certificate.cert.arn
}

resource "aws_cloudfront_distribution" "dist" {
  depends_on = [aws_acm_certificate_validation.validate]

  aliases = [local.domain_name]
  enabled = true
  price_class = "PriceClass_100"

  default_cache_behavior {
    allowed_methods = ["DELETE", "GET", "HEAD", "OPTIONS", "PATCH", "POST", "PUT"]
    cached_methods = ["GET", "HEAD"]
    min_ttl = 0
    default_ttl = 3600
    max_ttl = 86400
    target_origin_id = local.cfn_origin
    viewer_protocol_policy = "redirect-to-https"
    forwarded_values {
      query_string = false
      cookies {
        forward = "none"
      }
    }
  }
  origin {
    domain_name = replace(aws_apigatewayv2_api.api.api_endpoint, "/^https?://([^/]*).*/", "$1")
    origin_id = local.cfn_origin

    custom_origin_config {
      http_port = 80
      https_port = 443
      origin_protocol_policy = "https-only"
      origin_ssl_protocols = ["TLSv1.2"]
    }
  }
  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }
  viewer_certificate {
    acm_certificate_arn = aws_acm_certificate.cert.arn
    minimum_protocol_version = "TLSv1.2_2018"
    ssl_support_method = "sni-only"
  }
}

resource "aws_route53_record" "A" {
  name = local.domain_name
  type = "A"
  zone_id = data.aws_route53_zone.zone.zone_id

  alias {
    evaluate_target_health = false
    name = aws_cloudfront_distribution.dist.domain_name
    zone_id = aws_cloudfront_distribution.dist.hosted_zone_id
  }
}