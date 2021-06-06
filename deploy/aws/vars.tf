variable "aws_region" {
  type        = string
  description = "The AWS region to deploy into"
  default     = "us-east-2"
}

variable "auth0_domain" {
  type = string
  description = "Domain for Auth0 JWT authorizer"
  default = "dev-2usxzn4i.us.auth0.com"
}

variable "auth0_audience" {
  type = string
  description = "The audience for the authorizor"
  default = "https://auth0-jwt-authorizer"
}