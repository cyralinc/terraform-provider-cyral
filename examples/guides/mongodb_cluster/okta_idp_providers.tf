locals {
  # Replace [TENANT] by your tenant name. Ex: mycompany.app.cyral.com
  control_plane_host = "[TENANT].app.cyral.com"
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

# Refer to okta provider documentation:
#
# * https://registry.terraform.io/providers/okta/okta/latest/docs
#
provider "okta" {
  org_name  = "dev-123456" # your organization name
  base_url  = "okta.com"   # your organization url
  api_token = "xxxx"
}

provider "aws" {
  region = "us-east-1"
}
