---
page_title: "Provider: Cyral"
description: |-
  Use this provider to interact with resources supported by Cyral. You must provide proper credentials before you can use it.
---

# Provider

The provider is the base element and must be used to inform application-wide
parameters, like the Cyral control plane reference and authentication secrets.

## Example Usage

- Terraform v12

```hcl
provider "cyral" {
    client_id = ""     # optional
    client_secret = "" # optional
    control_plane = "some-cp.cyral.com:8000"
}
```

- Terraform v13+

```hcl
terraform {
  required_providers {
    cyral = {
      source = "cyralinc/cyral"
      version = ">= 2.2.0"
    }
  }
}

provider "cyral" {
    client_id = ""     # optional
    client_secret = "" # optional
    control_plane = "some-cp.cyral.com:8000"
}
```

## Argument Reference

* `auth_provider` - (Optional) Authorization provider in use by the Control Plane (valid values: `auth0`, `keycloak`). Default: `keycloak`.
* `auth0_audience` - (Optional) Auth0 audience.
* `auth0_domain` - (Optional) Auth0 domain name.
* `client_id` - (Optional) Client id used to authenticate against the Control Plane.
* `client_secret` - (Optional) Client secret used to authenticate against the Control Plane.
* `control_plane` - (Required) Control plane host and API port (ex: `some-cp.cyral.com:8000`)

----

Authentication parameters `client_id` and `client_secret` are defined as optional in the provider body once they can be set through environment variables in order to avoid storing secrets in source code repositories. The environment variables corresponds to `CYRAL_TF_CLIENT_ID` and `CYRAL_TF_CLIENT_SECRET` respectivelly and can be defined as follows:

- Linux/Mac

```bash
export CYRAL_TF_CLIENT_ID=""
export CYRAL_TF_CLIENT_SECRET=""
```

- Windows

```
set CYRAL_TF_CLIENT_ID=""
set CYRAL_TF_CLIENT_SECRET=""
```

### Provider Credentials - UI

#### New Credentials

A `Service Account` must be created in order to use the provider. It can be created through the control plane UI, accessing the `Service accounts` section in the left menu and clicking on the `+` button. Choose a name for the new service account and select the following roles so you can use all the provider functions:

<img src="https://raw.githubusercontent.com/cyralinc/terraform-provider-cyral/main/docs/images/create_service_account.png">

Confirm the account creation by clicking on the `CREATE` button. This will generate a `Client ID` and a `Client Secret` that should be used in the [provider configuration](#example-usage).

#### Rotate Credentials

To rotate secrets for existing service accounts, select a specific service account in the UI, and then click on the button `ROTATE CLIENT SECRET` as the image below suggests:

<img src="https://raw.githubusercontent.com/cyralinc/terraform-provider-cyral/main/docs/images/rotate_client_secret.png">

That will generate a new `Client Secret` that you can copy and use to replace the old one.

### Provider Credentials - Script

#### New Credentials

A `Service Account` must be created in order to use the provider. It can be created by the [script provided in the scripts folder](../scripts/create-keycloak-service-account.sh). You can run it with the command below:

```bash
curl https://raw.githubusercontent.com/cyralinc/terraform-provider-cyral/main/scripts/create-keycloak-service-account.sh -O
bash create-keycloak-service-account.sh
```

#### Rotate Credentials

[This script](../scripts/rotate-keycloak-service-account-secret.sh) can be used to rotate secrets for existing service accounts. It can be rotated by running the command below:

```bash
curl https://raw.githubusercontent.com/cyralinc/terraform-provider-cyral/main/scripts/rotate-keycloak-service-account-secret.sh -O
bash rotate-keycloak-service-account-secret.sh
```
