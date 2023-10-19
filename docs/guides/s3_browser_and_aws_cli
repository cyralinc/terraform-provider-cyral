---
page_title: "S3 File Browser and AWS CLI"
---

Use this guide to create the minimum required configuration in both Cyral
Control Plane and your AWS account to deploy a Cyral Sidecar to AWS EC2
to protect your S3 bucket and allow connections from [AWS CLI](https://cyral.com/docs/connect/s3-connect/cli)
and [Cyral S3 File Browser](https://cyral.com/docs/connect/s3-connect/s3-browser).

By running this example you will have a fully functional sidecar on your AWS
account. Read the comments and update the necessary parameters as instructed.

See also the [Cyral Sidecar module for AWS EC2](https://registry.terraform.io/modules/cyralinc/sidecar-ec2/aws/latest)
for more details on how the sidecar is deployed to AWS and more advanced configurations.

See also the [S3 File Browser](https://cyral.com/docs/manage-repositories/s3/s3-sidecar)
and [SSO for S3](https://cyral.com/docs/manage-repositories/s3/s3-sso) documentation
for more details about the feature.

```terraform
terraform {
  required_providers {
    cyral = {
      source  = "cyralinc/cyral"
      version = "~> 4.7"
    }
  }
}

locals {
  # Replace [TENANT] by your tenant name. Ex: mycompany.app.cyral.com
  control_plane_host = "[TENANT].app.cyral.com"

  # Use the name of the IdP that will be used to access the S3 Browser
  idp = {
    name = "<IDP_NAME_AS_SHOWN_IN_THE_UI>"
  }

  repos = {
    s3 = {
      # These are the ports the sidecar will accept connections
      # for S3 browser and S3 CLI
      browser_port = 443
      cli_port = 453
    }
  }

  sidecar = {
    # Set to true if you want a sidecar deployed with an
    # internet-facing load balancer (requires a public subnet).
    public_sidecar = true

    # Set the desired sidecar version or leave it empty if
    # you prefer to control the version from the control plane
    # (later only possible in CPs >=v4.10).
    sidecar_version = "v4.10.1"

    # Set the AWS region that the sidecar will be deployed to
    region = ""
    # Set the ID of VPC that the sidecar will be deployed to
    vpc_id = ""
    # Set the IDs of the subnets that the sidecar will be deployed to
    subnets = [""]
    # Name of the CloudWatch log group used to push logs
    cloudwatch_log_group_name = "cyral-example-loggroup"

    # Set the allowed CIDR block for SSH access to the sidecar
    ssh_inbound_cidr = ["0.0.0.0/0"]
    # Set the allowed CIDR block for database access through the
    # sidecar
    db_inbound_cidr = ["0.0.0.0/0"]
    # Set the allowed CIDR block for monitoring requests to the
    # sidecar
    monitoring_inbound_cidr = ["0.0.0.0/0"]

    # Set the ARN for the certificate that will be used by the load balancer
    # for S3 Browser connections
    load_balancer_certificate_arn = ""
    # Set the hosted zone ID that will be used to create the DNS name in
    # parameter `dns_name`
    dns_hosted_zone_id = ""
    # Set the DNS name that will be used by your sidecar. Ex:
    # sidecar.mycompany.com
    dns_name = ""
  }
}

provider "aws" {
  region = local.sidecar.region
}

# Follow the instructions in the Cyral Terraform Provider page to set
# up the credentials:
#
# * https://registry.terraform.io/providers/cyralinc/cyral/latest/docs
provider "cyral" {
  client_id     = ""
  client_secret = ""

  control_plane = local.control_plane_host
}

# The log group is created in AWS by module.cyral_sidecar
# when the sidecar is deployed.
resource "cyral_integration_logging" "cloudwatch" {
  name = "my-cloudwatch"
  cloudwatch {
    region = local.sidecar.region
    group  = local.sidecar.cloudwatch_log_group_name
    stream = "cyral-sidecar"
  }
}

resource "cyral_sidecar" "sidecar" {
  name               = "my-sidecar"
  deployment_method  = "terraform"
  activity_log_integration_id = cyral_integration_logging.cloudwatch.id
}

resource "cyral_sidecar_credentials" "sidecar_credentials" {
  sidecar_id = cyral_sidecar.sidecar.id
}

resource "cyral_repository" "s3" {
  name = "s3repo"
  type = "s3"

  repo_node {
    host = "s3.amazonaws.com"
    port = 443
  }
}

resource "cyral_sidecar_listener" "s3_cli" {
  sidecar_id = cyral_sidecar.sidecar.id
  repo_types = ["s3"]
  network_address {
    port = local.repos.s3.cli_port
  }
  s3_settings {
    proxy_mode = true
  }
}

resource "cyral_sidecar_listener" "s3_browser" {
  sidecar_id = cyral_sidecar.sidecar.id
  repo_types = ["s3"]
  network_address {
    port = local.repos.s3.browser_port
  }
  s3_settings {
    proxy_mode = false
  }
}

resource "cyral_repository_binding" "s3" {
  sidecar_id    = cyral_sidecar.sidecar.id
  repository_id = cyral_repository.s3.id
  listener_binding {
    listener_id = cyral_sidecar_listener.s3_cli.listener_id
  }
  listener_binding {
    listener_id = cyral_sidecar_listener.s3_browser.listener_id
  }
}

data "cyral_integration_idp_saml" "saml" {
    display_name = local.idp.name
}

# Let users from the provided `identity_provider` use SSO
# to access the database
resource "cyral_repository_conf_auth" "s3" {
  repository_id     = cyral_repository.s3.id
  identity_provider = data.cyral_integration_idp_saml.saml.idp_list[0].id
}

# Enables the access portal for this repository in the
# especified sidecar
resource "cyral_repository_access_gateway" "s3" {
  repository_id = cyral_repository.s3.id
  sidecar_id    = cyral_sidecar.sidecar.id
  binding_id    = cyral_repository_binding.s3.binding_id
}

###########################################################################
# Creates an IAM policy that the sidecar will assume in order to access
# your S3 bucket. In this example, the policy attached to the role will
# let the sidecar access all buckets.

data "aws_iam_policy_document" "s3_access_policy" {
  statement {
    actions   = ["s3:*"]
    resources = [
      "arn:aws:s3:::*"
    ]
  }
}

resource "aws_iam_policy" "s3_access_policy" {
  name        = "sidecar_s3_access_policy"
  path        = "/"
  description = "Allow sidecar to access S3"
  policy      = data.aws_iam_policy_document.s3_access_policy.json
}

data "aws_iam_policy_document" "sidecar_trust_policy" {
  statement {
    actions = ["sts:AssumeRole"]
    effect  = "Allow"
    principals {
      type        = "AWS"
      identifiers = [module.cyral_sidecar.aws_iam_role_arn]
    }
  }
}

resource "aws_iam_role" "s3_role" {
  name               = "sidecar_s3_access_role"
  path               = "/"
  assume_role_policy = data.aws_iam_policy_document.sidecar_trust_policy.json
}

resource "aws_iam_role_policy_attachment" "s3_role_policy_attachment" {
  role       = aws_iam_role.s3_role.name
  policy_arn = aws_iam_policy.s3_access_policy.arn
}
###########################################################################

resource "cyral_repository_user_account" "s3_repo_user_account" {
  name = aws_iam_role.s3_role.arn
  repository_id = cyral_repository.s3.id
  auth_scheme {
    aws_iam {
      role_arn = aws_iam_role.s3_role.arn
    }
  }
}

# Set the proper identity for the username, email or group that will
# be allowed to access the S3 browser
resource "cyral_repository_access_rules" "access_rule" {
  repository_id = cyral_repository.s3.id
  user_account_id = cyral_repository_user_account.s3_repo_user_account.user_account_id
  rule {
    identity {
      type = "email"
      name = "myuser@mycompany.com"
    }
  }
}

module "cyral_sidecar" {
  source = "cyralinc/sidecar-ec2/aws"

  # Use the module version that is compatible with your sidecar.
  version = "~> 4.3"

  sidecar_version = local.sidecar.sidecar_version

  sidecar_id = cyral_sidecar.sidecar.id
  control_plane = local.control_plane_host
  client_id     = cyral_sidecar_credentials.sidecar_credentials.client_id
  client_secret = cyral_sidecar_credentials.sidecar_credentials.client_secret

  sidecar_ports = [local.repos.s3.browser_port, local.repos.s3.cli_port]

  vpc_id  = local.sidecar.vpc_id
  subnets = local.sidecar.subnets

  ssh_inbound_cidr        = local.sidecar.ssh_inbound_cidr
  db_inbound_cidr         = local.sidecar.db_inbound_cidr
  monitoring_inbound_cidr = local.sidecar.monitoring_inbound_cidr

  load_balancer_scheme        = local.sidecar.public_sidecar ? "internet-facing" : "internal"
  associate_public_ip_address = local.sidecar.public_sidecar


  load_balancer_certificate_arn = local.sidecar.load_balancer_certificate_arn
  load_balancer_tls_ports       = [
    local.repos.s3.browser_port
  ]

  sidecar_dns_hosted_zone_id = local.sidecar.dns_hosted_zone_id
  sidecar_dns_name           = local.sidecar.dns_name
}

output "sidecar_load_balancer_dns" {
  value = module.cyral_sidecar.sidecar_load_balancer_dns
}
```
