output "gw_invoke_url" {
  value = aws_apigatewayv2_stage.v1.invoke_url
}

output "cfn_domain" {
  value = aws_cloudfront_distribution.dist.domain_name
}