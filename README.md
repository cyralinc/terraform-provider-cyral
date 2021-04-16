
# Cyral Terraform Provider

The Cyral Terraform Provider contains resources that can be used to interact with the Cyral API through Terraform code. It allows customers to maintain a history of changes in Cyral environment by transforming configurations into code and use Terraform to control state changes.

Our provider uses the same naming conventions and organization as stated in Terraform guidelines for [writing custom providers](https://www.terraform.io/docs/extend/writing-custom-providers.html).

## Change Log

The [Change Log](CHANGELOG.md) keeps track of features, bug fixes and Control Plane compatibility of this provider.

## Compatibility

This provider is compatible with both Auth0 or Keycloak-based CPs. Some initial setup is needed in both Auth0 and Keycloak as stated in the next sections.

### Auth0

1. Open Auth0 dashboard;
2. Select `Applications` and hit `Create application`;
    1. Choose `Machine to Machine Applications`;
    2. Select the API `https://cyral-api.com`;
    3. Select scopes: `read:users`;
    4. Finish the creation by clicking `Authorize`;
3. In the application just created, access `Settings` and copy `Client ID` and `Client Secret`. Use these parameters to set up the provider. See the [provider](./doc/provider.md) documentation how to set those two parameters.

### Keycloak

A `Service Account` for the Terraform provider must be created using the following steps:

```
# Define the target control plane
export CONTROL_PLANE=mycontrolplane.cyral.com

# Get a token from the CP (use the UI or some API)
export TOKEN="Authorization: Bearer ..."

# Get role ids necessary to run the provider
export ROLE_IDS=`curl -X GET https://$CONTROL_PLANE:8000/v1/users/roles -H "$TOKEN" -H "Content-type:Application/JSON" | jq '[.roles | map(select(.name | contains("Modify Integrations", "Modify Policies", "Modify Roles","Modify Sidecars and Repositories", "View Sidecars and Repositories"))) | .[].id]' -c`

# Create/update a service account for Terraform and return the necessary parameters
# client_id and client_secret that will be used in the provider
curl -X POST https://$CONTROL_PLANE:8000/v1/users/serviceAccounts -d '{"displayName":"terraform","roleIds":'"$ROLE_IDS"'}' -H "$TOKEN" -H "Content-type:Application/JSON" | jq

echo "Use both `client_id` and `client_secret` returned above to set up your provider."
echo "See the provider documentation in https://github.com/cyralinc/terraform-provider-cyral/blob/main/doc/provider.md how to set those two parameters."
echo
```

## Usage Example

The code below is just a simple example of how to use the Cyral Terraform Module. Refer to the "Supported Elements" section for more information on resources and provider details.

### Terraform v0.12

```hcl
provider "cyral" {
    client_id = ""     # optional
    client_secret = "" # optional
    control_plane = "some-cp.cyral.com:8000"
}

resource "cyral_repository" "my_repo_name" {
    host = "myrepo.cyral.com"
    port = 3306
    type = "mariadb"
    name = "myrepo"
}

resource "cyral_integration_elk" "elk" {
    name = "my-elk-integration"
    kibana_url = "kibana.local"
    es_url = "es.local"
}

resource "cyral_integration_datadog" "datadog" {
    name = "my-datadog-integration"
    api_key = "datadog-api-key"
}

resource "cyral_sidecar" "my_sidecar_name" {
    name = "mysidecar"
    deployment_method = "cloudFormation"
    log_integration = cyral_integration_elk.elk.id
    metrics_integration_id = cyral_integration_datadog.datadog.id
    aws_configuration {
        publicly_accessible = false
        aws_region = "us-east-1"
        key_name = "ec2-key-name"
        vpc = "vpc-id"
        subnets = "subnetid1,subnetid2,subnetidN"
    }
}

resource "cyral_datamap" "my_datamap_name" {
    mapping {
        label = "CNN"
        data_location {
            repo = cyral_repository.my_repo_name.name
            attributes = ["applications.customers.credit_card_number"]
        }
    }
}
```

### Terraform v0.13 and v0.14

```hcl
terraform {
  required_providers {
    cyral = {
      source = "cyral.com/terraform/cyral"
    }
  }
}

provider "cyral" {
    auth0_domain = "some-name.auth0.com"
    auth0_audience = "cyral-api.com"
    control_plane = "some-cp.cyral.com:8000"
}

resource "cyral_repository" "mongodb_repo" {
    type = "mongodb"
    host = "mongodb.cyral.com"
    port = 27017
    name = "mymongodb"
}

resource "cyral_repository" "mariadb_repo" {
    type = "mariadb"
    host = "mariadb.cyral.com"
    port = 3307
    name = "mymariadb"
}

resource "cyral_integration_elk" "elk" {
    name = "my-elk-integration"
    kibana_url = "kibana.local"
    es_url = "es.local"
}

resource "cyral_integration_datadog" "datadog" {
    name = "my-datadog-integration"
    api_key = "datadog-api-key"
}

resource "cyral_sidecar" "my_sidecar_name" {
    name = "mysidecar"
    deployment_method = "cloudFormation"
    log_integration = cyral_integration_elk.elk.id
    metrics_integration_id = cyral_integration_datadog.datadog.id
    aws_configuration {
        publicly_accessible = false
        aws_region = "us-east-1"
        key_name = "ec2-key-name"
        vpc = "vpc-id"
        subnets = "subnetid1,subnetid2,subnetidN"
    }
}

locals {
    repositories = [cyral_repository.mongodb_repo, cyral_repository.mariadb_repo]
}

resource "cyral_repository_binding" "repo_binding" {
    count         = length(local.repositories)
    repository_id = local.repositories[count.index].id
    listener_port = local.repositories[count.index].port
    sidecar_id    = cyral_sidecar.my_sidecar_name.id
}

resource "cyral_datamap" "my_datamap_name" {
    mapping {
        label = "CNN"
        data_location {
            repo = cyral_repository.my_repo_name.name
            attributes = ["applications.customers.credit_card_number"]
        }
    }
}
```

## Supported Operations

Terraform Cyral Provider was designed to be compatible with all existing Terraform operations. Thus it supports `apply`, `destroy`, `graph`, `import`, `refresh`, `show`, `taint`, `untaint`, etc.

### Import

Import operation identifies resources using property `name`. Thus, if you need to import the state of the resource `cyral_repository.my_resource_name` shown above, you will run:

```shell
terraform import cyral_repository.my_resource_name myrepo
```

## Supported Elements
- [Data Source Sidecar Template](./doc/data_source_sidecar_template.md)
- [Provider](./doc/provider.md)
- [Resource Datamap](./doc/resource_datamap.md)
- [Resource Identity Map](./doc/resource_identity_map.md)
- [Resource Integration Datadog](./doc/resource_integration_datadog.md)
- [Resource Integration ELK](./doc/resource_integration_elk.md)
- [Resource Integration Logstash](./doc/resource_integration_logstash.md)
- [Resource Integration Looker](./doc/resource_integration_looker.md)
- [Resource Integration Okta](./doc/resource_integration_okta.md)
- [Resource Integration Microsoft Teams](./doc/resource_integration_microsoft_teams.md)
- [Resource Integration Slack Alerts](./doc/resource_integration_slack_alerts.md)
- [Resource Integration Splunk](./doc/resource_integration_splunk.md)
- [Resource Integration Sumo Logic](./doc/resource_integration_sumo_logic.md)
- [Resource Policy](./doc/resource_policy.md)
- [Resource Policy Rule](./doc/resource_policy_rule.md)
- [Resource Repository](./doc/resource_repository.md)
- [Resource Repository Authentication Configuration](./doc/resource_repository_conf_auth.md)
- [Resource Repository Binding](./doc/resource_repository_binding.md)
- [Resource Repository Local Account](./doc/resource_repository_local_account.md)
- [Resource Sidecar](./doc/resource_sidecar.md)

## Prerequisites

Our existing provider supports Terraform `v0.12`, `v0.13` and `v0.14`. There are special actions to be taken in order to use this provider with Terraform `v0.12` as described in the `Deployment` section.

## Build Instructions

In order to build and distribute this provider, follow the steps below:

 1. Clone [terraform-provider-cyral](https://github.com/cyralinc/terraform-provider-cyral) repo from GitHub;

 2. Go to the root directory of the cloned repo using Linux shell and execute `make`. The build process will create binaries in directory `out` for both `darwin` and `linux` 64 bits. These binaries will be copied automatically to the local Terraform registry to be used by Terraform 13 and 14.

## Test Instructions

The test framework requires basic configuration before it can be executed as follows:

1. Set the configuration environment variables:

```bash
# Set the control plane DNS name and port (default 8000):
export CYRAL_TF_CP_URL=mycp.cyral.com:8000

# Set Keycloak client and secret ID:
export CYRAL_TF_CLIENT_ID=?
export CYRAL_TF_CLIENT_SECRET=?

# Initialize Terraform acceptance tests variable
export TF_ACC=true
```

2. Run `make`

## Deployment

### Terraform v0.12

Copy the desired binary file created in directory `out` (see "Build Instructions") to the root folder containing those `.tf` files that will be used to handle Cyral Terraform provider resources.

Run `terraform init` and proceed with `terraform apply` normally to execute your Terraform scripts.

### Terraform v0.13 and v0.14

If you **are** running the provider with the same user and machine you built the provider using steps in `Build Instructions`, you should just run `terraform init` and proceed with `terraform apply` normally to execute your Terraform scripts.

If you **are not** running the provider with the same user *or* are not in the same machine that you built the provider, you must copy the binaries in directory `out` to the local registry as follows:

```bash
cd terraform-provider-cyral
cp out/${OS_ARCH}/${BINARY} ~/.terraform.d/plugins/cyral.com/terraform/cyral/${VERSION}/${OS_ARCH}
```

Where:
* **OS_ARCH** corresponds to the distribution (`darwin_amd64` or `linux_amd64`);
* **BINARY** corresponds to the binary name. Ex: `terraform-provider-cyral_v0.1.0`;
* **VERSION** corresponds to the version number withouth `v`. Ex: `0.1.0`.
