# Configures `my-sidecar-cloudwatch` to push logs to CloudWatch to a log
# group named `cyral-example-loggroup` and a stream named `cyral-sidecar`.
locals {
  cloudwatch_log_group_name = "cyral-example-loggroup"
}

resource "cyral_sidecar" "sidecar" {
  name                        = "my-sidecar-cloudwatch"
  deployment_method           = "terraform"
  activity_log_integration_id = cyral_integration_logging.cloudwatch.id
}

resource "cyral_integration_logging" "cloudwatch" {
  name = "my-cloudwatch"
  cloudwatch {
    region = "us-east-1"
    group  = local.cloudwatch_log_group_name
    stream = "cyral-sidecar"
  }
}

resource "cyral_sidecar_credentials" "creds" {
  sidecar_id = cyral_sidecar.sidecar.id
}

module "cyral_sidecar" {
  source  = "cyralinc/sidecar-ec2/aws"
  version = "~> 4.0"

  sidecar_id = cyral_sidecar.sidecar.id

  cloudwatch_log_group_name = local.cloudwatch_log_group_name

  client_id     = cyral_sidecar_credentials.creds.client_id
  client_secret = cyral_sidecar_credentials.creds.client_secret

  # ...
}
