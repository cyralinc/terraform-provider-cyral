
# Cyral Terraform Provider

The Cyral Terraform Provider contains resources that can be used to interact with the Cyral API through Terraform code. It allows customers to maintain a history of changes in Cyral environment by transforming configurations into code and use Terraform to control state changes.

Our provider uses the same naming conventions and organization as stated in Terraform guidelines for [writing custom providers](https://www.terraform.io/docs/extend/writing-custom-providers.html).


## Documentation

Full and comprehensive documentation for this provider is available on the [provider documentation index](./docs/index.md).

### Change Log

The [Change Log](CHANGELOG.md) keeps track of features, bug fixes and Control Plane compatibility of this provider.

### Guides

See below a list of guides that can be used to deploy some predefined scenarios:

- [Add native repository credentials to AWS Secrets Manager](./docs/guides/native_credentials_aws_sm.md)


## Building and Testing

### Build Instructions

In order to build this repository, follow the steps below:

 1. Clone [terraform-provider-cyral](https://github.com/cyralinc/terraform-provider-cyral) repo from GitHub;

 2. Go to the root directory of the cloned repo using Linux shell and execute `make`. The build process will create binaries in directory `out` for both `darwin` and `linux` 64 bits. These binaries will be copied automatically to the local Terraform registry to be used by Terraform 13 and later.

Alternatively, you can use the dockerfile to build the image using `make docker-compose/build`

To use the local provider, the module must be configured to use the local provider path as follows:
```hcl
terraform {
  required_providers {
    cyral = {
      source = "local/terraform/cyral"
    }
  }
}
```

### Test Instructions

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

### Running project built locally

#### Terraform v0.12

Copy the desired binary file created in directory `out` (see [Build Instructions](build-instructions#)) to the root folder containing those `.tf` files that will be used to handle Cyral Terraform provider resources.

Run `terraform init` and proceed with `terraform apply` normally to execute your Terraform scripts.

### Terraform v0.13+

**If you are running** the provider with the same user and machine you built the provider using steps in [Build Instructions](build-instructions#), you should just run `terraform init` and proceed with `terraform apply` normally to execute your Terraform scripts.

**If you are not running** the provider with the same user *or* are not in the same machine that you built the provider, you must copy the binaries in directory `out` to the local registry as follows:

```bash
cd terraform-provider-cyral
cp out/${OS_ARCH}/${BINARY} ~/.terraform.d/plugins/cyral.com/terraform/cyral/${VERSION}/${OS_ARCH}
```

Where:
* **OS_ARCH** corresponds to the distribution (`darwin_amd64` or `linux_amd64`);
* **BINARY** corresponds to the binary name. Ex: `terraform-provider-cyral_v0.1.0`;
* **VERSION** corresponds to the version number withouth `v`. Ex: `0.1.0`.
