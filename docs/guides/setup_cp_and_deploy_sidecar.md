---
page_title: "Create a sidecar and repository in Cyral control plane and bring up the sidecar on AWS"
---

Use the following code to get a basic scenario with sidecar and repository set in the Cyral 
control plane and then bring this sidecar up on your AWS account.

By running this example you will have a fully functional sidecar on your AWS account to
protect your database. Read the comments and update the necessary parameters as instructed.

See also the [Cyral sidecar module for AWS](https://registry.terraform.io/modules/cyralinc/sidecar-aws/cyral/latest)
for more details on how the sidecar is deployed to AWS.

```hcl
terraform {
  required_providers {
    cyral = {
      source = "cyralinc/cyral"
      version = ">= 2.4.0"
    }
  }
}

locals {
    # Replace [TENANT] by your tenant name. Ex: mycompany.app.cyral.com
    control_plane = "[TENANT].app.cyral.com"

    sidecar = {
        # The sidecar name prefix is a unique name that is used by the sidecar module to create
        # resources in the target AWS account. For this reason, we use 'cyral-zzzzzz' where
        # zzzzzz are the last 6 digits of the sidecar id created in the control plane. This
        # explanation is purely informational and we advise you to keep this variable as is.
        sidecar_name_prefix = "cyral-${substr(lower(cyral_sidecar.mysql_sidecar.id), -6, -1)}"

        # If you would like to use other log integration, download a new template from the UI
        # and copy the log integration configuration or follow this module documentation.
        log_integration = "cloudwatch"
        # Set the desired EC2 instance type for the auto scaling group.
        instance_type = "t3.medium"
        # Set to true if you want a sidecar deployed with an internet-facing load
        # balancer (requires a public subnet).
        public_sidecar = false

        # Set the AWS region that the sidecar will be deployed to
        region = ""
        # Set the ID of VPC that the sidecar will be deployed to
        vpc_id = ""
        # Set the IDs of the subnets that the sidecar will be deployed to
        subnets = [""]

        # Set the allowed CIDR block for SSH access to the sidecar
        ssh_inbound_cidr = ["0.0.0.0/0"]
        # Set the allowed CIDR block for database access through the sidecar
        db_inbound_cidr = ["0.0.0.0/0"]
        # Set the allowed CIDR block for health check requests to the sidecar
        healthcheck_inbound_cidr = ["0.0.0.0/0"]

        # Set the parameters to access the private Cyral container registry.
        # These parameters can be found on the sidecar Terraform template
        # downloaded from the UI.
        container_registry = {
            name         = "" # see container_registry in the downloaded template
            username     = "" # see container_registry_username in the downloaded template
            registry_key = "" # see container_registry_key in the downloaded template
        }
    }
}

provider "aws" {
    region = local.sidecar.region
}

provider "cyral" {
    # Follow the instructions in the Cyral Terraform Provider page to set up the 
    # credentials: https://registry.terraform.io/providers/cyralinc/cyral/latest/docs
    client_id = ""
    client_secret = ""

    control_plane = "${local.control_plane}:8000"
}

resource "cyral_sidecar" "mysql_sidecar" {
    name = "MySQLSidecar"
    deployment_method = "terraform"
}

resource "cyral_sidecar_credentials" "sidecar_credentials" {
    sidecar_id = cyral_sidecar.mysql_sidecar.id
}

resource "cyral_repository" "mysql_repo" {
    name = "mysql_repo"
    type = "mysql"
    host = "mysql.mycompany.com"
    port = 3306
}

resource "cyral_repository_binding" "repo_binding" {
    repository_id = cyral_repository.mysql_repo.id
    listener_port = cyral_repository.mysql_repo.port
    sidecar_id    = cyral_sidecar.mysql_sidecar.id
}

module "cyral_sidecar" {
    # Set the desired sidecar version. This information can be extracted
    # from the template downloaded from the UI.
    sidecar_version = "v2.27.0"

    source  = "cyralinc/sidecar-aws/cyral"
    # Use the module version that is compatible with your sidecar. This
    # information can be extracted from the template downloaded from 
    # the UI.
    version = "2.5.4"

    sidecar_id = cyral_sidecar.mysql_sidecar.id 

    name_prefix = local.sidecar.sidecar_name_prefix

    control_plane = local.control_plane

    repositories_supported = ["mysql"]

    sidecar_ports = [cyral_repository.mysql_repo.port]

    mongodb_port_alloc_range_low  = 0
    mongodb_port_alloc_range_high = 0

    instance_type   = local.sidecar.instance_type
    log_integration = local.sidecar.log_integration
    vpc_id          = local.sidecar.vpc_id
    subnets         = local.sidecar.subnets

    ssh_inbound_cidr         = local.sidecar.ssh_inbound_cidr
    db_inbound_cidr          = local.sidecar.db_inbound_cidr
    healthcheck_inbound_cidr = local.sidecar.healthcheck_inbound_cidr

    load_balancer_scheme        = local.sidecar.public_sidecar ? "internet-facing" : "internal"
    associate_public_ip_address = local.sidecar.public_sidecar

    deploy_secrets   = true
    secrets_location = "/cyral/sidecars/${cyral_sidecar.mysql_sidecar.id}/secrets"

    container_registry          = local.sidecar.container_registry.name
    container_registry_username = local.sidecar.container_registry.username
    container_registry_key      = local.sidecar.container_registry.registry_key
    client_id                   = cyral_sidecar_credentials.sidecar_credentials.client_id
    client_secret               = cyral_sidecar_credentials.sidecar_credentials.client_secret
}

output "sidecar_dns" {
    value = module.cyral_sidecar.sidecar_dns
}

output "sidecar_load_balancer_dns" {
    value = module.cyral_sidecar.sidecar_load_balancer_dns
}
```
