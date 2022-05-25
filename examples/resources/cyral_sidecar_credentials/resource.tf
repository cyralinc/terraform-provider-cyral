terraform {
  required_providers {
    cyral = {
      source = "cyralinc/cyral"
    }
  }
  backend "s3" {
    bucket = "some-s3-bucket"
    key    = "terraform-state.json"
    region = "us-east-1"
    encrypt = true
    dynamodb_table = "some-dynamodb-table"
  }
}
