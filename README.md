
# Cyral Terraform Provider

The Cyral Terraform Provider contains resources that can be used to interact with the Cyral API through Terraform code. It allows customers to maintain a history of changes in Cyral environment by transforming configurations into code and use Terraform to control state changes.

Our provider uses the same naming conventions and organization as stated in Terraform guidelines for [writing custom providers](https://www.terraform.io/docs/extend/writing-custom-providers.html).

## Usage Example

The code below is just a simple example of how to use the Cyral Terraform Module. Refer to the "Supported Elements" section for more information on resources and provider details.

```hcl
provider "cyral" {
    auth0_domain = "some-name.auth0.com"
    auth0_audience = "cyral-api.com"
    control_plane = "some-cp.cyral.com:8000"
}

resource "cyral_repository" "my_resource_name" {
    host = "myrepo.cyral.com"
    port = 3306
    type = "mariadb"
    name = "myrepo"
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

- [Provider](./doc/provider.md)
- [Resource Repository](./doc/resource_repository.md)

## Prerequisites

Our existing provider supports Terraform `v0.12`, `v0.13` and `v0.14`. There are special actions to be taken in order to use this provider with Terraform `v0.12` as described in the `Deployment` section.

## Build Instructions

In order to build and distribute this provider, follow the steps below:

 1. Clone [terraform-provider-cyral](https://github.com/cyralinc/terraform-provider-cyral) repo from GitHub;

 2. Go to the root directory of the cloned repo using Linux shell and execute `make`. The build process will create binaries in directory `out` for both `darwin` and `linux` 64 bits. These binaries will be copied automatically to the local Terraform registry to be used by Terraform 13 and 14.


## Deployment

### Terraform v0.12

Copy the desired binary file created in directory `out` (see "Build Instructions") to the root folder containing those `.tf` files that will be used to handle Cyral Terraform provider resources.

Run `terraform init` and proceed with `terraform apply` normally to execute your Terraform scripts.

### Terraform v0.13

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