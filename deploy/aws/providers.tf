terraform {
  required_version = ">= 0.15"

  backend "remote" {
    organization = "bkimbrough"

    workspaces {
      name = "resume-backend"
    }
  }

  required_providers {
    archive = {
      source  = "hashicorp/archive"
      version = "~> 2.2.0"
    }
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.4"
    }
  }
}

provider "archive" {}

provider "aws" {
  region = var.aws_region
  default_tags {
    tags = {
      app       = "resume"
      component = "backend"
      managedBy = "Terraform"
    }
  }
}

provider "aws" {
  alias  = "usea1"
  region = "us-east-1"
  default_tags {
    tags = {
      app       = "resume"
      component = "backend"
      managedBy = "Terraform"
    }
  }
}