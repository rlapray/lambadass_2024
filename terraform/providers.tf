terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.47"
    }
  }
  required_version = ">= 1.7.4"
}

provider "aws" {
  region  = var.region
  profile = var.profile
  default_tags {
    tags = {
      Owner   = var.owner
      Project = var.project
    }
  }
}
