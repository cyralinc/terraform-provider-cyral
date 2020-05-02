
# Cyral Terraform Provider

The Cyral Terraform Provider contains resources that can be used to interact with the Cyral API through Terraform code. It allows customers to maintain a history of changes in Cyral environment by transforming configurations into code and use Terraform to control state changes.

Our provider uses the same naming conventions and organization as stated in Terraform guidelines for [writing custom providers](https://www.terraform.io/docs/extend/writing-custom-providers.html).

## Usage Example

The code below is just a simple example of how to use the Cyral Terraform Module. Refer to the "Supported Elements" section for more information on resources and provider details.

```hcl
provider "cyral" {
    auth0_domain = "some-name.auth0.com"
    auth0_audience = "cyral-api.com"
    control_plane = "some-cp.cyral.com"
}

resource "cyral_repository" "my_resource_name" {
    host = "myrepo.cyral.com"
    port = 3306
    type = "mariadb"
    name = "myrepo"
    require_tls = true
}
```

## Supported Elements

- [Provider](./doc/provider.md)
- [Resource Repository](./doc/resource_repository.md)

## Build Instructions

In order to build and distribute this provider, follow the steps below:

 1. Clone [terraform-provider-cyral](https://github.com/cyralinc/terraform-provider-cyral) repo from GitHub;

 2. Go to the root of the cloned repo using Linux shell and execute `make`. A binary file named `terraform-provider-cyral` must be created in the same directory. It corresponds to the provider that will be used to deploy to Terraform.

## Deployment

Copy the binary file created in the "Build Instructions" session to the root of the folder containing those `.tf` files that will be used in your Terraform session.

Next, run `terraform init` and proceed with `terraform apply` normally to execute your Terraform scripts.
